CSC 464 - Concurrency
Prof. Yvonne Coady
Assignment #1

# Introduction

Assignment #1 is a great opportunity for student's in Professor Yvone Coady's 
Concurrency class (UVIC CSC 464/564) to experience, measure, and attempt to solve 
concurrency problems. 
The assignment consists of 5 problems from the Little Book of Semaphores by 
Allen B. Downey (Version 2.2.1), and 1 problem that students have arrived at 
on their own. 
The problems from the Little Book of Semaphores are chosed by students, 
and intended to model real life problems and scenarios. 
The 6th problem (student's choice) is inteded to be an interesting problem, 
or one it's related to their final project.

This document is broken down into the following sections. 
This introduction section, outlining the assignment parameters. 
A methodology section, describing how I performed measurements and 
analysis of software implementations to problems as well as where I 
gathered inforamation about the tooling. 
The discussion section, where each problem is broken down into it's relevance, 
code/runtime characteristics, and analysis. 
And lastly, a conclusion section, which contains a highlevel overview of the 
project, realizations, and closing remarks.

# Methodology

For each problem, I wrote two or more implementations. 
These either implemented solutions in a different fashion in the same language, 
or a similar fashion in a different language. 

For solutions in GoLang, I use a tool called `pprof`. 
It allowed me to collect cpu and memory resource information about the go program. 
In particular, I can see exactly how long it took to execute, 
what routines took however long to execute, how much memory was used, 
and what lines of code were responsible. 
Resources used to understand this tool are included 
[here (understanding profiling tool)](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/), 
[here (profiling tool documentation & setup)](https://blog.golang.org/profiling-go-programs), 
[here (memory profiling gotcha)](https://austburn.me/blog/go-profile.html),
[and here (runtime pprof tool)](https://golang.org/pkg/runtime/pprof/)/ 
[here (http pprof tool)](https://golang.org/pkg/net/http/pprof/)

In Node, I'm using the built in V8 profiler. 
It breaks down what languages inside the program are used for what amount of time. 
Diving in futher, it breaks down what events use how much time. 
Since we've only done one problem in Node, and we're mostly interested in it's 
runtime in various scenarios, we are not diving into memory usage 
(also because it requires additional libraries). 
Resourced used to understand Node's V8 profiling tools are included 
[here](https://blog.ghaiklor.com/profiling-nodejs-applications-1609b77afe4e).

And in Python, I use their built-in profiling tools. 
In particular `profile` and `cProfile` for script time, function execution times, 
and cpu and memory resource consumption.
I found out and learned most from the python docs 
[found here](https://docs.python.org/3/library/profile.html).

After performing measurements on each problem, 
I analyse the different implementations by comparing and contrasting 
their correctness, comprehesability, and performance. 
For some problems some of these cases are harder to make an argument for. 
However, in each scenario I do my best to use logic and metrics, 
backing them up with logical pathways or reasonable justification
(which I hope is correct!).

# Discussion

## (1) Producer/Consumer

### Relevance

Producers and consumers are analogous to many systems in computing today. 
On a network level it is the server and client, exchanging information. 
On a local level it could be a peripheral communicating with the host. 
And on a distributed level it could be multiple sensors producing data which 
is being consumed and managed by a database system (or multiple!). 

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

Winner: Both


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

#### Correctness

#### Comprehesibility

#### Performance

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

#### Correctness

#### Comprehesibility

#### Performance

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

Let's compare the code the output commented out 
as we've previously shown it slows down the program significantly. 
Comparing the un-buffered channel vs the buffered channel implementation, 
we want to see if adding a buffer to the channels allows for faster processing. 
Since each O must wait for two Hs to tell it that they've bonded, and then the Hs 
must wait for the O to tell them they've bonded, the idea behind the buffered 
channel is that fewer context switches may occur. 
This is because the Os can tell Hs it's ready and wait for them to bond, 
and the Hs don't need to wait for the O to acknowledge they've bonded since 
there may only be two Hs in the critical section at once.

The un-buffered channel implementation's stats:

```
Duration: 2.16s, Total samples = 6.36s (295.06%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 10
Showing nodes accounting for 5.94s, 93.40% of 6.36s total
Dropped 46 nodes (cum <= 0.03s)
Showing top 10 nodes out of 41
      flat  flat%   sum%        cum   cum%
     4.81s 75.63% 75.63%      4.81s 75.63%  runtime.usleep
     0.41s  6.45% 82.08%      0.41s  6.45%  runtime.pthread_cond_signal
     0.20s  3.14% 85.22%      0.20s  3.14%  runtime.procyield
     0.16s  2.52% 87.74%      0.16s  2.52%  runtime.pthread_mutex_lock
     0.08s  1.26% 88.99%      0.89s 13.99%  runtime.chanrecv
     0.07s  1.10% 90.09%      0.08s  1.26%  runtime.gfget
     0.07s  1.10% 91.19%      5.10s 80.19%  runtime.lock
     0.06s  0.94% 92.14%      0.07s  1.10%  runtime.stackpoolalloc
     0.04s  0.63% 92.77%      0.04s  0.63%  runtime.casgstatus
     0.04s  0.63% 93.40%      0.11s  1.73%  runtime.malg.func1
```

```
Showing nodes accounting for 8.50MB, 100% of 8.50MB total
      flat  flat%   sum%        cum   cum%
    8.50MB   100%   100%     8.50MB   100%  runtime.malg
         0     0%   100%     8.50MB   100%  runtime.mstart
         0     0%   100%     8.50MB   100%  runtime.newproc.func1
         0     0%   100%     8.50MB   100%  runtime.newproc1
         0     0%   100%     8.50MB   100%  runtime.systemstack
```

And the buffered channel implementation's stats:

```
Duration: 2.05s, Total samples = 5.47s (266.97%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 10
Showing nodes accounting for 4.83s, 88.30% of 5.47s total
Dropped 39 nodes (cum <= 0.03s)
Showing top 10 nodes out of 50
      flat  flat%   sum%        cum   cum%
     3.64s 66.54% 66.54%      3.64s 66.54%  runtime.usleep
     0.33s  6.03% 72.58%      0.33s  6.03%  runtime.pthread_cond_signal
     0.27s  4.94% 77.51%      0.27s  4.94%  runtime.procyield
     0.13s  2.38% 79.89%      4.03s 73.67%  runtime.lock
     0.11s  2.01% 81.90%      0.13s  2.38%  runtime.gfget
     0.10s  1.83% 83.73%      0.10s  1.83%  runtime.casgstatus
     0.08s  1.46% 85.19%      0.08s  1.46%  runtime.stackpoolalloc
     0.07s  1.28% 86.47%      0.96s 17.55%  runtime.newproc1
     0.05s  0.91% 87.39%      0.07s  1.28%  runtime.gopark
     0.05s  0.91% 88.30%      0.05s  0.91%  runtime.pthread_mutex_lock
```

```
Showing nodes accounting for 8MB, 100% of 8MB total
      flat  flat%   sum%        cum   cum%
       8MB   100%   100%        8MB   100%  runtime.malg
         0     0%   100%        8MB   100%  runtime.mstart
         0     0%   100%        8MB   100%  runtime.newproc.func1
         0     0%   100%        8MB   100%  runtime.newproc1
         0     0%   100%        8MB   100%  runtime.systemstack
```

Seeing that this is suspiciously efficient. 
Let's try changing the code so it's as inneficient as possible and compare. 
To do this, we change: 

```
diff --git a/H2O/H2O-buffered.go b/H2O/H2O-buffered.go
index 9a4867e..30b2970 100644
--- a/H2O/H2O-buffered.go
+++ b/H2O/H2O-buffered.go
@@ -44,9 +44,12 @@ func main() {
 	}
 
 	for i:=0; i<1000000; i++ {
-		wg.Add(3)
-		go H()
+		wg.Add(1)
 		go O()
+	}
+	for i:=0; i<1000000; i++ {
+		wg.Add(2)
+		go H()
 		go H()
 	}
 	wg.Wait()
```

So that we create all the O routines, before creating any Hs. 

Now comparing the un-buffered vs. the buffered performance:

```
Duration: 9.27s, Total samples = 19.88s (214.38%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 10
Showing nodes accounting for 14100ms, 70.93% of 19880ms total
Dropped 84 nodes (cum <= 99.40ms)
Showing top 10 nodes out of 99
      flat  flat%   sum%        cum   cum%
    4030ms 20.27% 20.27%     4030ms 20.27%  runtime.usleep
    2720ms 13.68% 33.95%     2790ms 14.03%  runtime.stackpoolalloc
    1540ms  7.75% 41.70%     4030ms 20.27%  runtime.gentraceback
    1510ms  7.60% 49.30%     1640ms  8.25%  runtime.gcWriteBarrier
    1260ms  6.34% 55.63%     5260ms 26.46%  runtime.chanrecv
     800ms  4.02% 59.66%      860ms  4.33%  runtime.getempty
     620ms  3.12% 62.78%      950ms  4.78%  runtime.acquireSudog
     580ms  2.92% 65.69%     4780ms 24.04%  runtime.newproc1
     570ms  2.87% 68.56%      580ms  2.92%  runtime.gopark
     470ms  2.36% 70.93%     3260ms 16.40%  runtime.malg.func1
```

```
Showing nodes accounting for 765.22MB, 100% of 765.22MB total
      flat  flat%   sum%        cum   cum%
  748.27MB 97.79% 97.79%   748.27MB 97.79%  runtime.malg
   16.95MB  2.21%   100%    16.95MB  2.21%  runtime.allgadd
         0     0%   100%   765.22MB   100%  runtime.mstart
         0     0%   100%   765.22MB   100%  runtime.newproc.func1
         0     0%   100%   765.22MB   100%  runtime.newproc1
         0     0%   100%   765.22MB   100%  runtime.systemstack
```

vs. buffered:

```
Duration: 7.99s, Total samples = 19.23s (240.77%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 10
Showing nodes accounting for 13470ms, 70.05% of 19230ms total
Dropped 88 nodes (cum <= 96.15ms)
Showing top 10 nodes out of 102
      flat  flat%   sum%        cum   cum%
    5350ms 27.82% 27.82%     5360ms 27.87%  runtime.usleep
    2520ms 13.10% 40.93%     2540ms 13.21%  runtime.stackpoolalloc
     930ms  4.84% 45.76%      980ms  5.10%  runtime.getempty
     860ms  4.47% 50.23%      860ms  4.47%  runtime.pthread_cond_wait
     790ms  4.11% 54.34%      790ms  4.11%  runtime.pthread_cond_signal
     740ms  3.85% 58.19%     1430ms  7.44%  runtime.pcvalue
     680ms  3.54% 61.73%     1020ms  5.30%  runtime.gcWriteBarrier
     540ms  2.81% 64.53%     3630ms 18.88%  runtime.gentraceback
     540ms  2.81% 67.34%      880ms  4.58%  runtime.scanobject
     520ms  2.70% 70.05%      610ms  3.17%  runtime.step
```

```
Showing nodes accounting for 728.71MB, 99.93% of 729.21MB total
Dropped 1 node (cum <= 3.65MB)
      flat  flat%   sum%        cum   cum%
  711.76MB 97.61% 97.61%   711.76MB 97.61%  runtime.malg
   16.95MB  2.32% 99.93%    16.95MB  2.32%  runtime.allgadd
         0     0% 99.93%   728.71MB 99.93%  runtime.mstart
         0     0% 99.93%   728.71MB 99.93%  runtime.newproc.func1
         0     0% 99.93%   728.71MB 99.93%  runtime.newproc1
         0     0% 99.93%   728.71MB 99.93%  runtime.systemstack
```

Whoa, that's a difference!

2s and 8MB of memory usage, vs > 8s and over 700MB of memory usage. 

### Analysis

#### Correctness

The criteria for the H2O problem are:
1. H atoms must wait for an O atom in order to move forward
2. O atoms must wait for two H atoms in order to move forward

These are accomplished by synchronizing channels.

In the un-buffered implementation, the O routine holds until it's recieved two 
hReady signals from two separate H routines. 
At this point there are exactly two H routines waiting for the okay from exactly 
one O routine. 
And the O routine has now passed the barrier waiting for H routines. 
The O routine then release one O ready signal for one H routine to complete, and then one more for the other H routine to complete.

The difference with the buffered implementation is that it allows the Hs to enter 
the critical section without an O, but it does not let them through until they 
have been bound to by an O. 
Likewise, the Os may leave the crical section prior to Hs acknowledging it; 
but with the guarantee that there are exactly two Hs for each O, 
this is not a problem.

One caveat is that atoms passing are not necessarily the one allowed one another into the critical section. 
However, since we are still satisfying the criteria, this is not an issue. 

Winner: Both

#### Comprehesibility

Both solutions are short and simple. 
The un-buffered solution is slightly easier to mentally trace though, since the 
reader doesn't need to remember the size or what is in the buffered channel. 
They read through, see when a routine must wait, and when it sends a message. 
Otherwise one must conceptualize given the possible states of the buffered channel 
what coule happen.

Winner: Un-buffered

#### Performance

We see that the runtimes for 1000000 H2O moledules with an ideal execution path 
is only 2.16s and 2.05s, for the un-buffered and buffered channel 
implementations, respectively. 
And the memory usage is only 8.5MB and 8MB, respectively.
That means the difference in speed is less than 5%, and memory usage is only a 
0.5MB difference. 
From our experience comparing the producer/consumer problem, it's possible these 
differences are due to the runtime environment. 

However, when we look at the runtimes for 1000000 H2O molecules in a least 
optimal execution path, our runtimes jumps to 9.27s and 7.99s, respectively. 
And memory usage of 765.22MB and 729.21MB, respecitvely. 
What first seemed like a minor difference due to runtime variance, is 
exacerbated here given the un-ideal scenario. 
It's likely the extra routine switching in the un-buffered solution, 
causes the extra time and memory space. 

Winner: Buffered

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

#### Correctness

#### Comprehesibility

#### Performance

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

#### Correctness

#### Comprehesibility

#### Performance

# Conclusion