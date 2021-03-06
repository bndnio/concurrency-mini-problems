CSC 464 - Concurrency  
Prof. Yvonne Coady  
Assignment #1  

Written by: Brendon Earl

---

# Introduction

Assignment #1 is a great opportunity for students in Professor Yvone Coady's 
Concurrency class (UVIC CSC 464/564) to experience, measure, and attempt to solve 
concurrency problems. 
The assignment consists of 5 problems from the Little Book of Semaphores by 
Allen B. Downey (Version 2.2.1), and 1 problem that students have arrived at 
on their own. 
The problems from the Little Book of Semaphores are chosen by students, 
and intended to model real life problems and scenarios. 
The 6th problem (student's choice) is intended to be an interesting problem, 
or one it's related to their final project.

This document is broken down into the following sections. 
This introduction section outlines the assignment parameters. 
A methodology section, which describes how I measured solution implementations,
and where I gathered information about the tooling. 
The discussion section, where each problem is broken down into it's relevance, 
code/runtime characteristics, and analysis. 
And lastly, a conclusion section, which contains a high-level overview of the 
project, realizations, and closing remarks.

Code for each problem can be found in 
[my repository for these concurrency mini-problems](https://github.com/bndnio/concurrency-mini-problems).

# Methodology

For each problem, I wrote two or more implementations. 
These either implemented solutions in a different fashion in the same language, 
or a similar fashion in a different language. 
The performance is then evaluated using a performance tool for that language 
and comparing runtime and memory usage across the implementations with a fixed 
input size and similar execution order (e.g. H and O atom routines are called 
in a similar order). 

For implementations in GoLang, I use a tool called `pprof`. 
It allowed me to collect cpu and memory resource information about go programs. 
In particular, I can see exactly how long it took to execute, 
what routines took how long to execute, how much memory was used, 
and what lines of code were responsible. 
Resources used to understand this tool are included 
[here (understanding profiling tool)](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/), 
[here (profiling tool documentation & setup)](https://blog.golang.org/profiling-go-programs), 
[here (memory profiling gotcha)](https://austburn.me/blog/go-profile.html),
[and here (runtime pprof tool)](https://golang.org/pkg/runtime/pprof/)/ 
[here (http pprof tool)](https://golang.org/pkg/net/http/pprof/)

In Node, I'm using the built in V8 profiler. 
It breaks down what languages inside the program are used for what amount of time. 
Diving in further, it breaks down exactly what events use how much time. 
Since only solutions was written in Node, and we're mostly interested in it's 
runtime, we are not diving into memory usage 
(also because it requires additional libraries). 
Resourced used to understand Node's V8 profiling tools are included 
[here](https://blog.ghaiklor.com/profiling-nodejs-applications-1609b77afe4e).

And in Python, I use their built-in profiling tools. 
In particular `profile` and `cProfile` for script time, function execution times, 
and cpu and memory resource consumption.
I primarily learned from the python docs 
[found here](https://docs.python.org/3/library/profile.html).

After performing measurements on each problem, 
I analysed the different implementations by comparing and contrasting 
their correctness, comprehensibility, and performance. 
For some problems, some of these cases are harder to make an argument for. 
However, in each scenario I do my best to use logic and metrics, 
backing them up with reasonable justification (which I hope is correct!).

Note on the results: all measurements have been performed on my 
laptop which is not running in single user mode, or providing any contained environment. 
It is well understood that this isn't an ideal testing environment. 
However, for the purposes of this assignment, it is considered sufficient and 
relatively consistent (within the few minutes taken between tests of 
implementations for the same problem). 

# Discussion

## (1) Producer/Consumer

### Relevance

Producers and consumers are analogous to many systems in computing today. 
On a network level it is the server and client, exchanging information. 
On a local level it could be a peripheral communicating with the host. 
And on a distributed level it could be multiple sensors producing data which 
is being consumed and managed by a database system (or multiple!). 

### Code and Runtime Characteristics

The producer-consumer mechanism implemented consists of two separate routines. 
A single consumer is spun up to start. 
If it is successful in consuming data, it starts another go routine to help. 
If it can't consume data for two consecutive tries (spaced 100ms and 200ms apart), 
it becomes bored and kills itself. 

The main function spins up 100 producers, each of whom tries to enqueue 1000 
elements onto the shared data structure. 
This creates competition and simulates a scenario where multiple producers 
are sharing a data pool with multiple consumers.

Running `pprof` for the first time, and on prod-cons-channels, we see:

```
Duration: 10.30s, Total samples = 4.62s (44.84%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 3510ms, 75.97% of 4620ms total
Dropped 79 nodes (cum <= 23.10ms)
Showing top 10 nodes out of 108
      flat  flat%   sum%        cum   cum%
    1010ms 21.86% 21.86%     1010ms 21.86%  syscall.Syscall
     660ms 14.29% 36.15%      660ms 14.29%  runtime.pthread_cond_signal
     580ms 12.55% 48.70%      580ms 12.55%  runtime.pthread_cond_wait
     440ms  9.52% 58.23%      440ms  9.52%  runtime.usleep
     270ms  5.84% 64.07%      280ms  6.06%  runtime.stackpoolalloc
     150ms  3.25% 67.32%      510ms 11.04%  runtime.gentraceback
     120ms  2.60% 69.91%      120ms  2.60%  runtime.pthread_cond_timedwait_relative_np
     100ms  2.16% 72.08%      100ms  2.16%  runtime.nanotime
      90ms  1.95% 74.03%      110ms  2.38%  fmt.newPrinter
      90ms  1.95% 75.97%      200ms  4.33%  runtime.pcvalue
```

Noticing that `fmt.NewPrinter` is consuming a non-zero amount of time, 
let check that out a bit more.

```
(pprof) list main
Total: 4.62s
ROUTINE ======================== main.cons in /Users/brnd/repo/csc464/a1/prod-cons/prod-cons-channels.go
         0      1.84s (flat, cum) 39.83% of Total
         .          .     48:func cons() {
         .          .     49:   failedAttempts := 0
         .          .     50:   for {
         .          .     51:           wg.Add(1)
         .          .     52:
         .      310ms     53:           out, err := dequeue()
         .          .     54:           if (err == nil) {
         .          .     55:                   failedAttempts = 0
         .      910ms     56:                   fmt.Println("consumed: ", out)
         .          .     57:                   wg.Add(1)
         .       10ms     58:                   go cons()
         .          .     59:           } else {
         .          .     60:                   if (failedAttempts == 2) {
         .      540ms     61:                           fmt.Println("goodbye")
         .          .     62:                           break
         .          .     63:                   } else {
         .          .     64:                           failedAttempts++
         .       70ms     65:                           time.Sleep(time.Duration(100*failedAttempts)*time.Millisecond)
         .          .     66:                   }
         .          .     67:           }
         .          .     68:   }
         .          .     69:   wg.Done()
         .          .     70:}
...
```

Looks like our print statements takes 1.45s of the 1.84s 
it takes for the `cons` (consume) routine to run. 
This could be some sort of bottleneck! 
Let's remove it from the code base and rebuild and re-profile 
to find out.

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

Which looks very different, and takes way less time! 
Loosing the output to console I/O saved us nearly an order of magnitude on our runtime. 
It now looks like the most time is spend running usleep.
However, given that it's unknown if and when the next element will be produced, 
and the sleeping mechanism in the consumers allow us to prevent starving 
the producers, there isn't much left to do improve in this implementation.  

Before moving on, lets check the memory performance: 

```
(pprof) top
Showing nodes accounting for 94.15MB, 97.26% of 96.81MB total
Showing top 10 nodes out of 34
      flat  flat%   sum%        cum   cum%
   36.01MB 37.20% 37.20%    36.01MB 37.20%  runtime.malg
      20MB 20.66% 57.86%       20MB 20.66%  fmt.glob..func1
   10.28MB 10.61% 68.48%    10.28MB 10.61%  sync.(*Pool).Put
    6.50MB  6.71% 75.19%     6.50MB  6.71%  main.enqueue
    5.51MB  5.69% 80.88%     5.51MB  5.69%  time.Sleep
    5.35MB  5.52% 86.41%     5.35MB  5.52%  runtime.allgadd
    3.50MB  3.62% 90.03%    44.86MB 46.34%  runtime.systemstack
       3MB  3.10% 93.12%        3MB  3.10%  fmt.(*buffer).WriteByte
       2MB  2.07% 95.19%        2MB  2.07%  internal/poll.runtime_Semacquire
       2MB  2.07% 97.26%        2MB  2.07%  errors.New
```

Notice that it looks like we're using a total of 96.81MB.

Now, we can compare this implementation with channels to one using mutexes.

```
Duration: 1.14s, Total samples = 2.20s (192.32%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1460ms, 66.36% of 2200ms total
Dropped 38 nodes (cum <= 11ms)
Showing top 10 nodes out of 110
      flat  flat%   sum%        cum   cum%
     310ms 14.09% 14.09%      310ms 14.09%  runtime.usleep
     290ms 13.18% 27.27%      290ms 13.18%  runtime.pthread_cond_signal
     260ms 11.82% 39.09%      260ms 11.82%  runtime.pthread_cond_wait
     120ms  5.45% 44.55%      120ms  5.45%  runtime.memclrNoHeapPointers
     120ms  5.45% 50.00%      160ms  7.27%  runtime.stackpoolalloc
     100ms  4.55% 54.55%      100ms  4.55%  runtime.getempty
      90ms  4.09% 58.64%       90ms  4.09%  runtime.nanotime
      60ms  2.73% 61.36%      310ms 14.09%  runtime.gentraceback
      60ms  2.73% 64.09%       60ms  2.73%  runtime.pthread_cond_timedwait_relative_np
      50ms  2.27% 66.36%      120ms  5.45%  runtime.pcvalue
```

Looking at duration, it appears this is slightly slower! 
Though, wondering how much of that is due to runtime scheduling, let
run it again. 
This time, we get: 

```
Duration: 900.47ms, Total samples = 1.78s (197.67%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1180ms, 66.29% of 1780ms total
Showing top 10 nodes out of 132
      flat  flat%   sum%        cum   cum%
     300ms 16.85% 16.85%      300ms 16.85%  runtime.usleep
     190ms 10.67% 27.53%      200ms 11.24%  errors.New
     180ms 10.11% 37.64%      180ms 10.11%  runtime.pthread_cond_signal
     140ms  7.87% 45.51%      170ms  9.55%  runtime.stackpoolalloc
     120ms  6.74% 52.25%      140ms  7.87%  time.Sleep
      70ms  3.93% 56.18%       70ms  3.93%  runtime.pthread_cond_wait
      60ms  3.37% 59.55%       90ms  5.06%  runtime.acquireSudog
      40ms  2.25% 61.80%       40ms  2.25%  runtime.(*semaRoot).queue
      40ms  2.25% 64.04%      170ms  9.55%  runtime.gentraceback
      40ms  2.25% 66.29%       40ms  2.25%  runtime.getempty
```

Looks like it is now slightly faster than the channel implementation. 
Most of the performance difference must be due to how the routines are scheduled each run.

Checking the memory performance before moving on: 

```
(pprof) top
Showing nodes accounting for 71.22MB, 100% of 71.22MB total
Showing top 10 nodes out of 16
      flat  flat%   sum%        cum   cum%
   33.01MB 46.35% 46.35%    33.01MB 46.35%  runtime.malg
   14.69MB 20.63% 66.98%    14.69MB 20.63%  time.Sleep
    9.50MB 13.34% 80.32%     9.50MB 13.34%  sync.runtime_SemacquireMutex
    4.50MB  6.32% 86.64%     4.50MB  6.32%  errors.New (inline)
    4.33MB  6.08% 92.72%     4.33MB  6.08%  runtime.allgadd
    3.46MB  4.86% 97.58%    12.96MB 18.20%  main.enqueue
    1.72MB  2.42%   100%     1.72MB  2.42%  runtime/pprof.StartCPUProfile
         0     0%   100%    19.19MB 26.95%  main.cons
         0     0%   100%     4.50MB  6.32%  main.dequeue
         0     0%   100%     1.72MB  2.42%  main.main
```

Our memory use is a fair bit less than the channel implementation. 
Let's try running this one more time: 

```
(pprof) top
Showing nodes accounting for 78.02MB, 100% of 78.02MB total
Showing top 10 nodes out of 16
      flat  flat%   sum%        cum   cum%
   35.01MB 44.88% 44.88%    35.01MB 44.88%  runtime.malg
      13MB 16.66% 61.54%       13MB 16.66%  sync.runtime_SemacquireMutex
   10.86MB 13.92% 75.46%    10.86MB 13.92%  time.Sleep
    6.56MB  8.40% 83.87%    19.56MB 25.07%  main.enqueue
       6MB  7.69% 91.56%        6MB  7.69%  errors.New (inline)
    5.43MB  6.96% 98.52%     5.43MB  6.96%  runtime.allgadd
    1.16MB  1.48%   100%     1.16MB  1.48%  runtime/pprof.StartCPUProfile
         0     0%   100%    16.86MB 21.61%  main.cons
         0     0%   100%        6MB  7.69%  main.dequeue
         0     0%   100%     1.16MB  1.48%  main.main
```

This time memory usage is about 78MB. 
A bit strange seeing as it was only using about 71MB before. 
But it makes sense that memory usage can vary due to garbage collection. 
Also that our memory usage is lower than the channel version 
since there are no messages being passed around and being held in memory.

### Analysis

#### Correctness

Both implementations are a very easy to argue for correctness. 
The channel implementation is easy to argue since data pushed onto a buffered 
channel _should_ be safe. 
If there is too much data on the channel, that producer will simply wait. 
In the mutex implementation we see that each side can freely access the shared 
data structure, as long as they do so in a safe way using the mutex.

Both implementations however, do have the possibility of starvation. 
The channel implementation less-so since there is room to drop 100 elements into a 
buffered channel, giving time for the consumers to scale. 
Whereas in the mutex solution there is only one mutex every routine must fight 
to acquire. 
Nevertheless, if there are too many producers attempting to enqueue information, 
it could starve one that has been waiting a long time. 
The hope though, in both implementations is that the consumers scale with demand, 
and exit when bored. 
At least creating a scenario where producers can only be starved by other producers. 

One solution to this would be creating a queue for producers to sit in once they arrive. 
But this would encounter the same issue if too many hit the service at once. 
Though it would mitigate it since some order could be enforced, but this relies on 
scheduling what routine hits the queueing mechanism first, which isn't my responsibility.

Winner: Both


#### Comprehensibility

The comprehensibility is very similar in both cases. 
They both use the same skeleton, and exactly the same code that doesn't involve 
the locking mechanism. 

The difference though, is this:  

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

The channel solution is definitely more elegant. 
It's shorter, cleaner, no extra code controlling access. 
Just a buffered channel that producers can enqueue data to, and consumers 
can dequeue from.

Winner: Channel implementation

#### Performance

Given then runtime characteristics above, we see the biggest improvement of 
performance was omitting I/O in the form of printing to console. 
This brought the channel implementation runtime down from 9.25s to 1.07s! 
Comparing to the mutex implementation, runtime was relatively similar. 
The first run of the mutex implementation clocked 1.09s, while the second measured 1.06s. 
Laying on either side of the time of the channel implementation, it's fair so say 
the runtime is relatively equivalent, and depends more on scheduling at runtime.

The biggest difference between the two was the memory usage. 
Although the memory usage measurement was inconsistent, the usage was more than a 
20% reduction both times using mutexes instead of channels.

Winner: Mutex implementation

## (2) Readers/Writers

### Relevance

Like the producer/consumer and insert-search-delete problems, 
this is a problem which can be easily tied to a real application. 
Readers/writers is similar to producers/consumers in the way that there are 
two asymmetric actors interacting at a single point, being the data structure 
that they share. 
However, when writing, we need to ensure all readers are locked out and not 
active in order to retain the integrity of the data structure. 
This is similar to any software system which receives data which must be stored, 
that also fulfils requests to view this data. 
A particular example of this could be an RESTful http server. 
It could be receiving and executing requests simultaneously to write data 
to a data store and read from it. 
Where threads handling these requests need to control who has control to perform 
reads and writes. 

### Code and Runtime Characteristics

In this section we compare a readers/writers go program, 
with two variations of it written in Node.js. 
The first of the Node implementations is written with async and Promise await, 
while the second is done in a completely sequential nature. 
The reason for exploring this, is because Node runs on an event loop and we're not 
handling any I/O, there's not actually any asynchronous actions occurring. 
And when asynchronous function are queued up, then can they execute one at a time. 

First we check the time it takes for the Node scripts to process 1000000 
elements using bash's `time` library: 
we find that the async version runs according to: 
 
```
node readers-writers-async.js  1.71s user 0.27s system 213% cpu 0.925 total
```

and the sequential version's run stats are:  

```
node readers-writers-seq.js  0.11s user 0.02s system 95% cpu 0.133 total
```

What a huge difference! 
The sequential version is 15 times faster, and uses less than half the 
the computing power. 

Diving in, using Nodes `prof` and `prof-process` we see for 
the async implementation:

```
 [Summary]:
   ticks  total  nonlib   name
    150   20.5%   20.8%  JavaScript
    558   76.3%   77.3%  C++
    354   48.4%   49.0%  GC
      9    1.2%          Shared libraries
     14    1.9%          Unaccounted
```

vs. the sequential implementation

```
 [Summary]:
   ticks  total  nonlib   name
     12    3.4%    3.4%  JavaScript
    328   91.9%   93.7%  C++
      5    1.4%    1.4%  GC
      7    2.0%          Shared libraries
     10    2.8%          Unaccounted
```

We can see in the sequential process we're running has far fewer tics in all categories. 
Therefore meaning far less time was spent overall, and in every type of routine. 
Because we're not actually doing anything async, adding the async syntax 
likely just add significant overhead.

Now to see what our go program's cpu and memory usage looks like. 
Starting with 1000000 writers as we did with the Node scripts.

```
Duration: 625.19ms, Total samples = 1.28s (204.74%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 920ms, 71.88% of 1280ms total
Showing top 10 nodes out of 88
      flat  flat%   sum%        cum   cum%
     490ms 38.28% 38.28%      490ms 38.28%  runtime.usleep
     110ms  8.59% 46.88%      190ms 14.84%  sync.(*Mutex).Lock
      60ms  4.69% 51.56%       60ms  4.69%  runtime.memclrNoHeapPointers
      50ms  3.91% 55.47%       50ms  3.91%  runtime.procyield
      40ms  3.12% 58.59%      370ms 28.91%  main.reader
      40ms  3.12% 61.72%       40ms  3.12%  runtime.gfget
      40ms  3.12% 64.84%       40ms  3.12%  runtime.pthread_cond_wait
      30ms  2.34% 67.19%       30ms  2.34%  runtime.deferreturn
      30ms  2.34% 69.53%      440ms 34.38%  runtime.goexit0
      30ms  2.34% 71.88%       30ms  2.34%  runtime.pthread_cond_signal```

```

```
(pprof) top
Showing nodes accounting for 13429.24kB, 100% of 13429.24kB total
Showing top 10 nodes out of 13
      flat  flat%   sum%        cum   cum%
 6036.98kB 44.95% 44.95%  6036.98kB 44.95%  main.writer
 4097.50kB 30.51% 75.47%  4097.50kB 30.51%  runtime.malg
 1536.14kB 11.44% 86.90%  1536.14kB 11.44%  sync.runtime_SemacquireMutex
 1184.27kB  8.82% 95.72%  1184.27kB  8.82%  runtime/pprof.StartCPUProfile
  574.34kB  4.28%   100%   574.34kB  4.28%  runtime.allgadd
         0     0%   100%  1184.27kB  8.82%  main.main
         0     0%   100%  1536.14kB 11.44%  main.reader
         0     0%   100%  1184.27kB  8.82%  runtime.main
         0     0%   100%  4671.84kB 34.79%  runtime.mstart
         0     0%   100%  4671.84kB 34.79%  runtime.newproc.func1
```

and bumping up the input size to 10000000, 

```
Duration: 4.52s, Total samples = 11.85s (262.03%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 8650ms, 73.00% of 11850ms total
Dropped 92 nodes (cum <= 59.25ms)
Showing top 10 nodes out of 87
      flat  flat%   sum%        cum   cum%
    5110ms 43.12% 43.12%     5110ms 43.12%  runtime.usleep
    1100ms  9.28% 52.41%     1760ms 14.85%  sync.(*Mutex).Lock
     390ms  3.29% 55.70%      390ms  3.29%  runtime.pthread_cond_signal
     370ms  3.12% 58.82%      420ms  3.54%  runtime.gfget
     370ms  3.12% 61.94%      790ms  6.67%  sync.(*Mutex).Unlock
     360ms  3.04% 64.98%     3180ms 26.84%  main.reader
     340ms  2.87% 67.85%      340ms  2.87%  runtime.procyield
     210ms  1.77% 69.62%     2420ms 20.42%  runtime.newproc1
     200ms  1.69% 71.31%      200ms  1.69%  runtime.casgstatus
     200ms  1.69% 73.00%      260ms  2.19%  sync.(*WaitGroup).Add
```

```
(pprof) top
Showing nodes accounting for 63.20MB, 100% of 63.20MB total
Showing top 10 nodes out of 19
      flat  flat%   sum%        cum   cum%
   43.34MB 68.57% 68.57%    43.84MB 69.36%  main.writer
       8MB 12.66% 81.23%        8MB 12.66%  runtime.malg
       6MB  9.49% 90.73%        6MB  9.49%  sync.runtime_SemacquireMutex
    3.50MB  5.54% 96.26%    12.71MB 20.11%  runtime.systemstack
    1.20MB  1.91% 98.17%     1.20MB  1.91%  runtime.allgadd
    1.16MB  1.83%   100%     1.16MB  1.83%  runtime/pprof.StartCPUProfile
         0     0%   100%     1.16MB  1.83%  main.main
         0     0%   100%     5.50MB  8.70%  main.reader
         0     0%   100%     0.50MB  0.79%  math/rand.(*Rand).Int31
         0     0%   100%     0.50MB  0.79%  math/rand.(*Rand).Int31n
```

Looks like for 10 times the input size, our go program takes about 10 times the 
time and 4 times the memory space. 

But, our processing time is far greater than the sequential Node process! 
Thinking that GoLang is the fastest, especially in the face of a _scripting language_, 
would lead one to forget about some important overhead. 
Our go program has ONE MILLION routines running at once. 
Even if they are cheap, they're not free! 
Node, while maybe slower, runs everything sequentially, not worrying about 
coordinating, or sharing data, it just does.

### Analysis

#### Correctness

Evaluating correctness in node (for both the asynchronous and sequential implementation) 
is trivial. 
They are both absolutely correct, since only a single function may access the data 
regardless of if they're reading or writing. 

For the go implementation, it's a bit trickier. 
Here, a light-switch mechanism is applied. 
If a reader is the first in the room, it turns the light on 
(saying the room is not empty). 
Readers may come and go as they please from then until 
the last reader is going to leave. 
If a reader is the last to leave, they turn the light-switch off 
(mark room as empty). 
Once the room is empty, the writer may enter, and claim the room 
saying it is not empty (by turning the light-switch back on). 
At this point, any writer that attempts to enter must wait for the room 
to unlock again. 
And the first reader to attempt to enter will also be barred, while 
holding the readerCount Mutex. 
This ensures that no other reader can attempt to enter. 

This implementation does have the problem of starvation for writers. 
If the system is flooded with readers or many writers, then some writers 
may never get access to write their data. 

Winner: Node (in a cheat-y sort of way)

#### Comprehensibility

Due to the simplicity of the sequential implementation, 
Node is clearly more comprehensible. 
Without locks, or any concurrent behaviour to worry about it's 
far more comprehensible. 
Even the asynchronous version (which is not really async) 
is easy to see and understand. 

The go implementation requires the reader to have a stronger 
grasp of the state of shared variables, who is in the room, and 
who is allowed together. 
While the go implementation follows a well defined pattern, 
it isn't dead simple.

Winner: Node (in a cheat-y sort of way, again)

#### Performance

Depending on how you constrain the system, performance could go 
either way. 
In Node, it is reliable since there is no starvation and processes 
requests are handled in a first in first out (FIFO) manner. 

Given a large system to scale across, it's likely that go could 
perform better. 
But given the test environment on my laptop, Node ran faster. 

Winner: Node (in a cheat-y sort of way, again again)

## (3) Building H20

### Relevance

The H2O problem may initially seem strange, until one realizes it is an analogy 
for asynchronous systems which must synchronize. 
The two H atoms waiting to join the O atom, and then proceeding at the same time 
represents asynchronous functions which wait on each-other before exchanging 
data to complete the task together. 
An example of this could show up in data pipelines. 
Where at a stage multiple asynchronous handlers are dispatched to retrieve some 
information. 
Once they've completed, they must acknowledge one-another to return and continue 
the next stage of the pipeline.

### Code and Runtime Characteristics

Let's compare the code with print statements commented out, 
since we've previously seen that it significantly impacts performance. 
Comparing the un-buffered channel vs the buffered channel implementation, 
we want to see if adding a buffer to the channels allows for faster processing. 

Because each O must wait for two Hs to tell it that they've bonded, and then the Hs 
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
Let's try changing the code so it's as inefficient as possible and compare. 
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
Duration: 10.29s, Total samples = 22.91s (222.66%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 14050ms, 61.33% of 22910ms total
Dropped 110 nodes (cum <= 114.55ms)
Showing top 10 nodes out of 111
      flat  flat%   sum%        cum   cum%
    4380ms 19.12% 19.12%     4390ms 19.16%  runtime.usleep
    2480ms 10.82% 29.94%     2540ms 11.09%  runtime.stackpoolalloc
    1300ms  5.67% 35.62%     6160ms 26.89%  runtime.chanrecv
    1110ms  4.85% 40.46%     1140ms  4.98%  runtime.getempty
     950ms  4.15% 44.61%     4370ms 19.07%  runtime.gentraceback
     850ms  3.71% 48.32%      890ms  3.88%  runtime.gopark
     770ms  3.36% 51.68%      770ms  3.36%  runtime.pthread_cond_signal
     740ms  3.23% 54.91%     1240ms  5.41%  runtime.gcWriteBarrier
     740ms  3.23% 58.14%     1490ms  6.50%  runtime.scanobject
     730ms  3.19% 61.33%      800ms  3.49%  runtime.step
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

The criteria for the H2O problem is:  
1. H atoms must wait for an O atom in order to move forward  
2. O atoms must wait for two H atoms in order to move forward  

These are accomplished by synchronizing channels.

In the un-buffered implementation, the O routine holds until it's received two 
hReady signals from two separate H routines. 
At this point there are exactly two H routines waiting for the okay from exactly 
one O routine. 
Now the O routine has passed the barrier waiting for H routines. 
The O routine then release one oReady signal for one H routine to complete, 
and then one more for the other H routine to complete.

The difference with the buffered implementation is that it allows the Hs to enter 
the critical section without an O, but it does not let them through until they 
have been bound to by an O. 
Likewise, an O may leave the critical section prior it's Hs acknowledging it; 
but with the guarantee that there are exactly two Hs for each O, 
this is not a problem.

One caveat is that atoms passing are not necessarily the one that 
allowed one another into the critical section. 
However, since we are still satisfying the criteria, this is not an issue. 

Winner: Both

#### Comprehensibility

Both solutions are short and simple. 
The un-buffered solution is slightly easier to mentally trace though, since the 
reader doesn't need to remember the size or what is in the buffered channel. 
They read through, see when a routine must wait, and when it sends a message. 
Otherwise one must conceptualize the possible states of the buffered channel and
what could happen in each scenario.

Winner: Un-buffered

#### Performance

We see that the runtimes for 1000000 H2O molecules with an ideal execution path 
is only 2.16s and 2.05s for the un-buffered and buffered channel 
implementations, respectively. 
And the memory usage is only 8.5MB and 8MB, respectively.
That means the difference in speed is less than 5%, and memory usage is only a 
0.5MB difference. 
From our experience comparing the producer/consumer problem, it's possible these 
differences are due to the runtime environment. 

However, when we look at the runtimes for the same number of H2O molecules in a 
least optimal execution path, our runtimes jump to 9.27s and 7.99s, respectively. 
And memory usage to 765.22MB and 729.21MB, respectively. 
What first seemed like a minor difference due to runtime variance, is 
exacerbated here given the un-ideal scenario. 
It's likely the extra routine switching in the un-buffered solution, 
causes the extra time and memory space. 

Winner: Buffered

## (4) Insert-Search-Delete

### Relevance

The insert-search-delete is a relatable problem with many real-world applications. 
The first and most obvious, is for a data structure which can handle concurrency. 
This could be considered fundamental to implementation of concurrency at any scale, 
since without managing data, what else would happen with the output of a 
concurrent process. 
Less obviously, this could be synonymous for the various roles which actors 
(processes) would have on a concurrent system. 
Similar to the first point, this is simply a larger scale. 
For example, this state could be a user management system. 
Where employees can search the directory, managers take on the role of searching 
or inserting users, while administrators can take on the roll of the above 
or deleting users. 
Each of these actions should be done in a particular fashion in order 
to prevent data inconsistencies. 

Inserting searching and deleting are key components to any system, large or small. 
Handling these requests concurrently is a fundamental part in computing everywhere 
from local to distributed. 

### Code and Runtime Characteristics

The insert search delete code was mostly inspired by the 
"Little Book of Semaphores" pseudo code, and was noted as so. 

Here we do a direct comparison between using mutexes as 
implemented in GoLang's `sync` library, 
and using buffered channels to emulate mutexes. 

Running both with 100000 operators, with every second routine
being an inserter, every third (if not even) a searcher, 
and the remainder being deleters. 

Running, we found channels WAY slower. 
It took the channel implementation 5.89 minutes! 
While the mutexes only took 9.33 seconds! 

Taking a look at the performance output we see:

For the channel implementation.
```
Duration: 5.89mins, Total samples = 5.80mins (98.43%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 336.86s, 96.83% of 347.89s total
Dropped 178 nodes (cum <= 1.74s)
Showing top 10 nodes out of 54
      flat  flat%   sum%        cum   cum%
   183.75s 52.82% 52.82%    183.75s 52.82%  runtime.pthread_cond_signal
    97.88s 28.14% 80.95%    103.63s 29.79%  main.find
    37.58s 10.80% 91.76%     37.60s 10.81%  container/list.(*List).remove
     5.79s  1.66% 93.42%      5.79s  1.66%  container/list.(*Element).Next (inline)
     3.90s  1.12% 94.54%      3.91s  1.12%  runtime.usleep
     3.33s  0.96% 95.50%      3.33s  0.96%  runtime.pthread_cond_wait
     2.03s  0.58% 96.08%      2.05s  0.59%  runtime.stackpoolalloc
     1.41s  0.41% 96.49%      1.97s  0.57%  runtime.scanobject
     0.63s  0.18% 96.67%      3.92s  1.13%  runtime.newproc1
     0.56s  0.16% 96.83%      3.33s  0.96%  runtime.gentraceback
```

And the mutex implementation.  
```
Duration: 9.33s, Total samples = 19.97s (214.13%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 10890ms, 54.53% of 19970ms total
Dropped 112 nodes (cum <= 99.85ms)
Showing top 10 nodes out of 118
      flat  flat%   sum%        cum   cum%
    1690ms  8.46%  8.46%     1690ms  8.46%  runtime.usleep
    1650ms  8.26% 16.73%     1680ms  8.41%  runtime.stackpoolalloc
    1440ms  7.21% 23.94%     1660ms  8.31%  runtime.newdefer
    1060ms  5.31% 29.24%     1290ms  6.46%  main.find
     950ms  4.76% 34.00%     3870ms 19.38%  runtime.newproc1
     920ms  4.61% 38.61%     1370ms  6.86%  runtime.gcWriteBarrier
     900ms  4.51% 43.11%     1350ms  6.76%  runtime.acquireSudog
     810ms  4.06% 47.17%      810ms  4.06%  runtime.memclrNoHeapPointers
     750ms  3.76% 50.93%      750ms  3.76%  runtime.pthread_cond_signal
     720ms  3.61% 54.53%      750ms  3.76%  runtime.getempty
```

`runtime.pthread_cond_signal` was a huge resource hog for the 
channels. 

Diving in a little deeper into what lines took how long, 
we can see the find operation being super consuming.  

```
      10ms   1.72mins     65:   toRemove := find(value)
```

Where it didn't impact the mutex operation much at all. 
```
         .      360ms     65:   toRemove := find(value)
```

In the find operation, this is the impacting line.  
```
  1.58mins   1.58mins     43:           if e.Value == value {
```  
```
     1.05s      1.05s     43:           if e.Value == value {
```

One last thing to check for now is the memory usage, 
perhaps that will give us some clues.

The channel implementation used 54.8MBs, 
while the mutex implementation used 35.7MB. 
This isn't shocking, as we've seen above in the 
producer/consumer problem the memory characteristics 
followed a similar pattern.

To see if something was wrong with my implementation, 
I tried a sample size of 10000. 
Decreasing the number of routines by an order of magnitude. 
Somehow dropped the times down to 11 and 7 seconds!  
The drop to 7 seconds isn't unreasonable for the mutex 
implementation, but for the channels implementation 
this was revealing! 
The speedup was from almost 6 minutes to 11 seconds, for 
only a 10x change in input size. 
There must be some exponential running time aspect that 
I am not aware of.

It's strange how poorly the channel implementation behaved. 
And even stranger that the line with the biggest time impact 
had nothing to do with mutexes or channels. 
I hope that this is due to an anti-pattern I unknowningly followed, 
and not a poor characteristic of the language!

### Analysis

#### Correctness

Since the implementation using mutexes are channels are 
synonymous (only replacing lines with equivalent ones), 
we only need to argue for a single version that it's correct.

Each function is nice and small, so let's break it down 
one by one: 

```golang
func search(value int) *list.Element {
	defer wg.Done()
	searchSwitch.Lock(noSearcher)
	found := find(value)
	searchSwitch.Unlock(noSearcher)
	
	return found
}
```

Before searching, we `Lock` `searchSwitch`, meaning 
that if it's the first routine of it's type, 
it will take a lock on the parameter there (`noSearcher`). 
As long as it locks `noSearcher`, nothing which needs to 
bar searchers from operating can run. 
It then finds the value. 
Once found, it calls `Unlock` on `searchSearch`, 
which would unlock the mutex if it is the last routine there. 
And lastly, it returns the found value.

Next for `insert`, we have:  

```golang
func insert(value int) {
	defer wg.Done()
	insertSwitch.Lock(noInserter)
	insertMutex.Lock()
	l.PushBack(value)
	insertMutex.Unlock()
	insertSwitch.Unlock(noInserter)
}
```

It checks the 'lightswitch' `insertSwitch` on the way in, 
locking `noInserter` if need be. 
And then holds a mutex on the actual list to insert into. 
After inserting, it unlocks the mutex, and the runs 
`Unlock` on the `insertSwitch` on it's way out. 
This pattern protects the list when being appended to, 
and ensures no other code will run as long as it needs to 
bar insert routines.

Lastly, `delete` looks like this: 

```golang
func delete(value int) {
	defer wg.Done()
	noSearcher.Lock()
	noInserter.Lock()
	toRemove := find(value)
	if toRemove != nil { l.Remove(toRemove) }
	noInserter.Unlock()
	noSearcher.Unlock()
}
```

Being the most complicated of them all, delete needs to be 
since it performs the most impactful change. 
It locks out searchers and inserters (or waits to do so) 
with `noSearcher.Lock()` and `noInserter.Lock()`. 
Finds the value to remove with the find function, 
then removes it as long as that value was found. 
It then `Unlock`s `noSearcher` and `noInserter`.

Given the requirements:  
- many searchers can run at once  
- only a single inserter may run at once with many searchers concurrently  
- only a single delete routine may run with no concurrent searchers or inserters  

This code walk-through suggests that the implementation is correct.

Winner: Both

#### Comprehensibility

While I had thought the channels would end up being more 
comprehensible since they're the 'golden child' of GoLang. 
I believe using mutexes are actually much more readable. 

Comparing the lightswitch code next to each-other:

```golang
func (ls *lightswitch) Lock(m *sync.Mutex) {
	ls.mutex.Lock()
	ls.counter += 1
	if ls.counter == 1 {
		m.Lock()
	}
	ls.mutex.Unlock()
}
```

```golang
func (ls *lightswitch) Lock(c chan bool) {
	<- ls.mutex
	ls.counter += 1
	if ls.counter == 1 {
		<- c
	}
	ls.mutex <- true
}
```

It's clearer what is happening with the mutex version, 
in part because it's a common programming paradigm. 
We know that mutexes can only be held by a single 
routine, so it's obvious what is happening with 
`ls.mutex.Lock()` vs `<- ls.mutex`.
The details surrounding how a channel has been set up 
also isn't immediately clear, and so the developer 
working must remember that while tracing code. 

Winner: Mutex

#### Performance

Given the huge difference in the code and runtime 
characteristics section, the mutex implementation way 
outperforms the channel implementation. 
This is could be due to garbage collection of the channels, 
overhead being more prohibitive with large numbers of routines,
or the extra memory required.

Winner: Mutex

## (5) Sushi Bar

### Relevance

The sushi bar problem is similar to that of a type of controlled buffer or 
workload manager. 
Patrons represent jobs arriving somewhere for computation (at the sushi bar). 
Jobs can arrive and complete as necessary. 
But if the available space for jobs is completely used up, some cool-down mechanism 
is employed. 
In this case, the cool-down mechanism requires all jobs to finish before allowing more. 
In real-world scenarios, it may require job loads reduce to 50%. 
Justification could be for competing job requesters using shared resources. 
If a single requester is using all 100% of the resources, they are throttled 
to make room and prevent starvation for other requesters. 
In this example 'some other requester' could be a groups of system processes, 
garbage collection, or other peripherals/external systems.

### Code and Runtime Characteristics + Analysis

Skipped :(

## (6) Tangle Verification

### Relevance

The tangle is a concept in distributed computing. 
It is understood to be one of the competitors to blockchain due to it's lightweight 
nature. 
The tangle consists of the idea that in order for some work to be pushed onto the 
network, that actor must first do work to contribute to the network. 
It is a directed acyclic graph, where each piece of work is represented as a node. 
If an actor wishes to add some work to the system, it must verify the work 
of two others (nodes). 
Once other work is verified, a directed arrow is created from the nodes which 
received the verification to the node requesting work. 
Each node needs only two verifications, before it's considered 'verified' to 
the rest of the network.

For now, we will simplify the system, ignoring weighting of nodes to work on 
and instead treating it as a FIFO queue. 
An actor can read a node at the same time as any others, but to verify it, it 
must write the verified pointer to itself, so it must gain a write lock. 
Once an actor has verified two nodes and pointed to itself, it can add it's 
node (piece of work) to the queue for other actors to verify.

### Code and Runtime Characteristics

The Tangle verification, and tangle node code is strange. 
I had though prior to implementing it would be much simpler. 
The main issue I ran into, was managing the peer-to-peer 
communication where nodes are verifying one-another,
while having a supervisor manage a queue of nodes which 
still required verification and communicating this to nodes 
performing or needing work. 

Passing messages back and forth was complicated and 
the made the code obscure. 
Furthermore, issues would arise at random times when the 
count of routines rose above 50, 
but don't occur consistently under 90. 

I also wrote a python implementation which mirrors the 
implementation style of the go code by using queues 
instead of channels. 
However, once this block arose with the go code, 
I decided not to parallelize the python code for a few reasons:  
1. Because it was implemented with the same pattern, there would be no comprehensibility advantage,  
2. There is no way to benchmark it against the go code due to lack of stability  
3. It's likely this approach needs re-architecting, which defeats the purpose of the direct comparison  

## Analysis

### Correctness

The way of coordinating tasks among nodes and to the 
supervisor was hard to understand. 
The combination of the possibility of multiple concurrent 
jobs executing to verify a single node, while restricting 
the number of nodes which could claim the verification, 
and managing the queue of nodes to process made for some 
nasty spaghetti code. 

While the verification controller looked okay:
```golang
	defer wg.Done()
	newNode := tangle_node.New()

	var isComplete int
	for verified, i := 0, 0; verified < 2; i++ {
		var comm = make(chan bool, 1)
		var cb = make(chan bool)
		go queue[i%len(queue)].Verify(newNode, comm, cb)
		// ** node evaluation work would go here **
		comm <- true
		didVerify := <- cb
		if didVerify { verified++ }
	}
	<- queueWriteMutex
	queue = append(queue, newNode)
	if isComplete {
		queue = append(queue[:i%len(queue)], queue[((i%len(queue))+1):]...)
	}
      queueWriteMutex <- true
```

the tangle node got ugly quickly:

```golang
func (nd *Node) Verify(addr *Node, approved <-chan bool, cb chan<- bool) {
	// check if Node already approved
	<- nd.approvalMutex
	if nd.approvalCount >= 2 || inSlice(addr, nd.outbound) == true {
		fmt.Println("fast fail")
		fmt.Println(nd.approvalCount >= 2)
		nd.approvalMutex <- true
		cb <- false
		return
	}
	nd.approvalMutex <- true

	// allow requester to do work
	nd.workingMutex <- true
	var verified = <- approved
	if verified == true {
		// check for approval account again
		<- nd.approvalMutex
		// if full callback false
		// otherwise increase approval count
		if nd.approvalCount >= 2 || inSlice(addr, nd.outbound) == true  {
			fmt.Println("mid fail")
			cb <- false
			nd.approvalMutex <- true 
			<- nd.workingMutex
			return
		} else { 
			nd.approvalCount++
		}
		// and mark verified if possible
		if nd.approvalCount == 2 {
			nd.verified = true
		}
		// and set outbound link to Node
		nd.outbound = append(nd.outbound, addr)
		nd.approvalMutex <- true 
	}
	<- nd.workingMutex
	cb <- true
}
```

Winner: No-One

### Correctness

Because of the complexity, it's near impossible to 
construct a strong argument for correctness using static analysis. 
And because of the multiple execution paths that provide 
equivalent but not identical output, evaluating output for 
correctness is not simple either.

Because of the unknown looping that randomly occurred 
between 50 and 100 routines, it's fair to say this solution 
is probably not correct. 

And finally, because the python version followed the same 
logic and structure, it likely is not correct either.

Winner: No-One

### Performance

Getting stuck in an infinite loop without a clear reason is 
definitely not a performant feature.

To achieve decent performance, the implementation would 
have to be re-architectured so that it is working first. 
I should have started by creating nodes that only allowed a single 
verifier to work on it in order to reduce concurrent complexities. 
Then, it would have been appropriate to drive it _faster_ 
with more complex concurrent mechanisms and routines.

From this and and the insert-search-delete problem, it's 
made it apparent that mass concurrency and channels are 
not always the best solution. 
A well educated approach needs to be architected, and 
re-architected until a working solution is found. 
At that point it is safer to work forwards with more 
performant solutions. 

Winner: No-One

# Conclusion

This assignment has been a great experience in designing and writing, troubleshooting, and profiling concurrent software. 

A few great realizations that occurred from my work were that full blown concurrency is really not always faster.  
- Node sequential programming can be super efficient, and overhead can be a real problem in many languages (particularly Node and GoLang).  
- The number of routines in golang implementations can have a severe impact on performance.  
- The idea that more is not always faster really rang true on the insert-search-delete problem.  
- And the complexity introduced to early in the tangle problem made for a 'tangled' solution!  

Concurrent software is fun to think about and develop. 
But without care, it can be a real pain in the butt.
