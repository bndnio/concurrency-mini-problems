CSC 464 - Concurrency
Prof. Yvonne Coady
Assignment #1

# Introduction

# Discussion

## (1) Insert-Search-Delete

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

## (2) Producer/Consumer

### Relevance

Producers and consumers are analogous to many systems in computing today. 
On a network level it is the server and client, exchanging information. 
On a local level it could be a peripheral communicating with the host. 
And on a distributed level it could be multiple sensors producing data which 
is being consumed and managed by a database system. 

### Code and Runtime Characteristics

### Analysis

## (3) Readers/Writers

### Relevance

Like the previous two problems, this is a problem which can be easily tied to 
a real application. 
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