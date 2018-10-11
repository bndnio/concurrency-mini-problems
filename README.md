CSC 464 - Concurrency
Prof. Yvonne Coady
Assignment #1

# Introduction

# Discussion

## (1) Producer/Consumer

### Relevance

Producers and consumers are analogous to many systems in computing today. 
On a network level it is the server and client, exchanging information. 
On a local level it could be a peripheral communicating with the host. 
And on a distributed level it could be multiple sensors producing data which 
is being consumed and managed by a database system. 

### Code and Runtime Characteristics

Running pprof for the first time, we see:

```
Duration: 9.25s, Total samples = 4.27s (46.16%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 3140ms, 73.54% of 4270ms total
Dropped 81 nodes (cum <= 21.35ms)
Showing top 10 nodes out of 108
      flat  flat%   sum%        cum   cum%
     950ms 22.25% 22.25%      970ms 22.72%  syscall.Syscall
     480ms 11.24% 33.49%      480ms 11.24%  runtime.pthread_cond_signal
     480ms 11.24% 44.73%      480ms 11.24%  runtime.usleep
     460ms 10.77% 55.50%      460ms 10.77%  runtime.pthread_cond_wait
     260ms  6.09% 61.59%      290ms  6.79%  runtime.stackpoolalloc
     120ms  2.81% 64.40%      120ms  2.81%  fmt.newPrinter
     120ms  2.81% 67.21%      320ms  7.49%  main.dequeue
     120ms  2.81% 70.02%      480ms 11.24%  runtime.gentraceback
      80ms  1.87% 71.90%       90ms  2.11%  runtime.step
      70ms  1.64% 73.54%      190ms  4.45%  runtime.pcvalue
```

Noticing that `fmt.NewPrinter` is consuming a non-zero amount of time, 
let's remove it from the code base and rebuild and re-profile.

```
diff --git a/prod-cons-channels.go b/prod-cons-channels.go
index 7258bea..7db3f19 100644
--- a/prod-cons-channels.go
+++ b/prod-cons-channels.go
@@ -1,7 +1,6 @@
 package main
 
 import (
-	"fmt"
 	"sync"
 	"errors"
 	"time"
@@ -50,15 +49,15 @@ func cons() {
 	for {
 		wg.Add(1)
 
-		out, err := dequeue()
+		_, err := dequeue()
 		if (err == nil) {
 			failedAttempts = 0
-			fmt.Println("consumed: ", out)
+			// fmt.Println("consumed: ", out)
 			wg.Add(1)
 			go cons()
 		} else {
 			if (failedAttempts == 2) {
-				fmt.Println("goodbye")
+				// fmt.Println("goodbye")
 				break
 			} else {
 				failedAttempts++
```

Re-profiling gives us:

```
Duration: 1.07s, Total samples = 1.98s (184.46%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 1420ms, 71.72% of 1980ms total
Showing top 10 nodes out of 111
      flat  flat%   sum%        cum   cum%
     500ms 25.25% 25.25%      500ms 25.25%  runtime.usleep
     300ms 15.15% 40.40%      320ms 16.16%  runtime.stackpoolalloc
     150ms  7.58% 47.98%      150ms  7.58%  runtime.pthread_cond_signal
     100ms  5.05% 53.03%      100ms  5.05%  runtime.pthread_cond_wait
      90ms  4.55% 57.58%      200ms 10.10%  time.Sleep
      80ms  4.04% 61.62%      100ms  5.05%  runtime.scanobject
      50ms  2.53% 64.14%      260ms 13.13%  runtime.gentraceback
      50ms  2.53% 66.67%      730ms 36.87%  runtime.malg.func1
      50ms  2.53% 69.19%       50ms  2.53%  runtime.memclrNoHeapPointers
      50ms  2.53% 71.72%       80ms  4.04%  runtime.pcvalue
```

Which looks vastly different! And takes way less time! 
Loosing the output to console I/O saved us nearly an order of magnitude on our runtime. 
It now looks like sleeping is costing us the biggest performance hit. 
However, given that it's unknown if and when the next element will be produced, 
and the sleeping mechanism in the consumers allow us to prevent starving 
the producers, there isn't much left to do improve in this implementation.  

Before moving on, lets check the memory performance: 

```
Showing nodes accounting for 57.44MB, 100% of 57.44MB total
      flat  flat%   sum%        cum   cum%
   54.02MB 94.05% 94.05%    54.02MB 94.05%  runtime.malg
    1.89MB  3.29% 97.34%     1.89MB  3.29%  time.Sleep
    1.53MB  2.66%   100%     1.53MB  2.66%  runtime.allgadd
         0     0%   100%     1.89MB  3.29%  main.cons
         0     0%   100%    55.55MB 96.71%  runtime.mstart
         0     0%   100%    55.55MB 96.71%  runtime.newproc.func1
         0     0%   100%    55.55MB 96.71%  runtime.newproc1
         0     0%   100%    55.55MB 96.71%  runtime.systemstack
```

Noting that it looks like we're using a total of 57.44MB.

Now, we can compare this implementation with channels to one using mutexes!

```
Duration: 1.09s, Total samples = 1.93s (177.85%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 10
Showing nodes accounting for 1320ms, 68.39% of 1930ms total
Showing top 10 nodes out of 135
      flat  flat%   sum%        cum   cum%
     360ms 18.65% 18.65%      360ms 18.65%  runtime.usleep
     330ms 17.10% 35.75%      330ms 17.10%  runtime.pthread_cond_signal
     140ms  7.25% 43.01%      140ms  7.25%  runtime.pthread_cond_wait
     110ms  5.70% 48.70%      130ms  6.74%  time.Sleep
     100ms  5.18% 53.89%      140ms  7.25%  runtime.stackpoolalloc
      70ms  3.63% 57.51%       70ms  3.63%  runtime.(*semaRoot).queue
      60ms  3.11% 60.62%       80ms  4.15%  errors.New
      60ms  3.11% 63.73%      300ms 15.54%  runtime.gentraceback
      50ms  2.59% 66.32%      110ms  5.70%  runtime.pcvalue
      40ms  2.07% 68.39%       40ms  2.07%  runtime.nanotime
```

Looking at duration, it appears this is slightly slower! 
Though, wondering how much of that is due to runtime scheduling, let
running it again. 
This time, we get: 

```
Duration: 1.06s, Total samples = 2.04s (192.14%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 10
Showing nodes accounting for 1450ms, 71.08% of 2040ms total
Dropped 34 nodes (cum <= 10.20ms)
Showing top 10 nodes out of 95
      flat  flat%   sum%        cum   cum%
     430ms 21.08% 21.08%      430ms 21.08%  runtime.usleep
     290ms 14.22% 35.29%      290ms 14.22%  runtime.pthread_cond_signal
     160ms  7.84% 43.14%      160ms  7.84%  runtime.pthread_cond_wait
     110ms  5.39% 48.53%      110ms  5.39%  runtime.memclrNoHeapPointers
     110ms  5.39% 53.92%      130ms  6.37%  runtime.stackpoolalloc
      80ms  3.92% 57.84%      120ms  5.88%  errors.New
      80ms  3.92% 61.76%      110ms  5.39%  runtime.acquireSudog
      70ms  3.43% 65.20%      290ms 14.22%  runtime.gentraceback
      60ms  2.94% 68.14%      460ms 22.55%  runtime.newproc1
      60ms  2.94% 71.08%       60ms  2.94%  runtime.pthread_cond_timedwait_relative_np
```

Looks like it is now slightly faster than the channel implementation. 
Most of the performance difference must be due to how the routines are scheduled each run.

Checking the memory performance before moving on: 

```
Showing nodes accounting for 25.01MB, 100% of 25.01MB total
      flat  flat%   sum%        cum   cum%
   22.01MB 87.99% 87.99%    22.01MB 87.99%  runtime.malg
    1.91MB  7.64% 95.62%     1.91MB  7.64%  time.Sleep
    1.10MB  4.38%   100%     1.10MB  4.38%  runtime.allgadd
         0     0%   100%     1.91MB  7.64%  main.cons
         0     0%   100%    23.10MB 92.36%  runtime.mstart
         0     0%   100%    23.10MB 92.36%  runtime.newproc.func1
         0     0%   100%    23.10MB 92.36%  runtime.newproc1
         0     0%   100%    23.10MB 92.36%  runtime.systemstack
```

This is rather surprising! 
Our memory use is half that of the version using channels. 
Let's try running this one more time: 

```
Showing nodes accounting for 39249.53kB, 100% of 39249.53kB total
      flat  flat%   sum%        cum   cum%
36365.31kB 92.65% 92.65% 36365.31kB 92.65%  runtime.malg
 1121.44kB  2.86% 95.51%  1121.44kB  2.86%  main.enqueue
 1121.44kB  2.86% 98.37%  1121.44kB  2.86%  runtime.allgadd
  641.34kB  1.63%   100%   641.34kB  1.63%  time.Sleep
         0     0%   100%   641.34kB  1.63%  main.cons
         0     0%   100% 37486.75kB 95.51%  runtime.mstart
         0     0%   100% 37486.75kB 95.51%  runtime.newproc.func1
         0     0%   100% 37486.75kB 95.51%  runtime.newproc1
         0     0%   100% 37486.75kB 95.51%  runtime.systemstack
```

This time memory usage is almost 40Mb! 
A rather strange result seeing as it was only using about 25MB before. 

### Analysis

#### Correctness

Both implementations are a very easy to argue for correctness. 
The channel implementation is easy to argue since data pushed onto a buffered 
channel should always be safe. 
If there is too much data on the channel, that producer will simply wait. 
In the mutex implementation. 

Both implementations however, do have the posibility of starvation. 
The channel implementation less-so since there is room to drop 100 elements into a 
buffered channel, giving time for the consumers to scale. 
Whereas in the mutex solution there is only one mutex every routine must fight 
to acquire. 
Nevertheless, if there are too many producers attempting to enqueue information, 
it could starve one that has been waiting a long time. 
The hope though, in both implementations is that the consumers scale with demand, 
and exit when bored, At least creating a scenario where producers can only be 
starved by other producers. 

One solution to this would be creating a queue, but this would encounter the same 
issue if too many hit the service at once. 
Though it would mitigate it since some order could be enforced, and scheduling 
what routine hits the queueing mechanism first is the runtime mechanism's 
responsibility.


#### Comprehensibility

The comprehensibility is very similar in both cases. 
They both use the same skeleton, and exactly the same code that doesn't involve 
the locking mechanism. 

The difference though is this particularily:  

```golang
func enqueue(num int) {
    defer wg.Done()

	// *** START CRITICAL SECTION ***
	queue <- num
    // *** END CRITICAL SECTION ***
}
```

vs.

```golang
func enqueue(num int) {
    mutex.Lock()
	defer wg.Done()
	defer mutex.Unlock()

	// *** START CRITICAL SECTION ***
	queue = append(queue, num)
    // *** END CRITICAL SECTION ***
}
```

It's clear the channel solution is much more elegant. 
It's shorter, cleaner, no extra code controlling access. 
Just a buffered channel that producers can enqueue data to, and consumers 
can dequeue from.

Winner: Channel implementation

#### Performance

Give then analysis above, we see the biggest improvement of performance was 
omitting I/O in the form of printing to console. 
This brought the channel implementation runtime down from 9.25s to 1.07s! 
Comparing to the mutex implementation, runtime was relatively similar. 
The first run of the mutex implementation clocked 1.09s, while the second measured 1.06s. 
Laying on either side of the time of the channel implementation, it's fair so say 
the rumtime is relatively equivalent, and depends more on scheduling at runtime.

The biggest difference between the two was the memory usage. 
Although the memory usage measurement was inconsistent, the usage was more than a 
20% reduction both times using mutexes instead of channels.

Winner: Mutex implementation

## (2) Readers/Writers

### Relevance

Like the producer/consumer and insert-search-delete problems, 
this is a problem which can be easily tied to a real application. 
Readers/writers is similar to producers/consumers in the way that there are 
two asymetric actors interacting at a single point, being the data structure 
that they share. 
However here we need not worry about removing (or consuming) data, instead only 
reading it. 
This is similar to any software system which recieves data which must be stored, 
as well as requests to view this data. 
A particular example of this could be an RESTful http server. 
It could be recieving and executing requests simultaneously to write data 
to a data store and read from it. 

### Code and Runtime Characteristics

### Analysis

## (3) Insert-Search-Delete

### Relevance

The insert-serach-delete is a relatable problem with many real-world applications. 
The first and most obvious, is for a data structure which can handle concurrency. 
This could be considered fundamental to implementation of concurrency at any scale, since without managing data, what else would happen with the output of a 
concurrent process. 
Less obviously, this could be synonymous for the various roles which actors 
(processes) would have on a concurrent system. 
Similar to the first point, this is simply a larger scale. 
For example, this state could be a user management system. 
Where users can search of the directory, managers take on the role of searching 
or inserting users, while adminstrators can take on the roll of the above 
or deleting users. 
Each of these actions should be done in a mutual exclusive fashion in order 
to prevent data inconsistencies. And insertions (depending on the method) and 
deletion require even more care in regards to how many can influence the system 
at once. 

Inserting searching and deleting are key components to any system, large or small. 
Handling these requests concurrently is a fundamental part in computing everywhere 
from local to distributed. 

### Code and Runtime Characteristics

### Analysis

## (4) Building H20

### Relevance

The H2O problem may intially seem strange, until one realizes it is an anlogy 
for asynchronous systems which must synchronize. 
The two H atoms waiting to join the O atom, and then proceeding at the same time 
represents asynchronous functions which wait on eachother before exchanging 
data to complete the task together. 
An example of this could show up in data pipelines. 
Where at a stage multiple asynchronous handlers are dispatched to retrieve some 
information. 
Once they've completed, they must acknowledge one-another to return and continue 
the next stage of the pipeline.

### Code and Runtime Characteristics

### Analysis

## (5) Sushi Bar

### Relevance

The sushi bar problem is similar to that of a type of controlled buffer or 
workload manager. 
Patrons represent jobs arriving somewhere for computation (at the sushi bar). 
Jobs can arrive and complete as necessary. 
But if the available space for jobs is completely used up, some cool-down mechanism 
is employed. 
In the case the cool-down mechanism requires all jobs to finish before allowing more. 
In real-world scenarios, it may require job loads reduce to 50%. 
Justification could be for competing job sources using shared resources. 
If a single requester is using all 100% of the resources, they are throttled 
to make room and prevent starvation for other requesters. 
Some other requester could be for critical system processes, garbage collection, 
or other peripherals/external systems.

### Code and Runtime Characteristics

### Analysis

## (6) Tangle Verification

### Relevance

The tangle is a concept in distributed computing. 
It is understood to be one of the competitors to blockchain due to it's lightweight 
nature. 
The tangle consists of the idea that in order for some work to be pushed onto the 
network, it must first do work to contribute to the network. 
It is a directed acyclic graph, where each piece of work is represented as a node. 
As an actor wishes to add some work to the system, it must verify the work 
of two others (nodes). 
When verified, a directed arrow is created from the nodes which recieved the 
verification to the node requesting work. 

For now, we will simplify the system, ignoring weighting and instead treating 
work as a FIFO queue. 
An actor can read a node at the same time as any others, but to verify it, it 
must write the verified pointer to itself, so it must gain a write lock. 
Once an actor has verified two nodes and pointed to itself, it can add it's 
node (piece of work) to the queue for other actors to verify.

### Code and Runtime Characteristics

### Analysis

# Conclusion