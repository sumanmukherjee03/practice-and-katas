### k-way merge

This is a helpful pattern to solve problems that involve a list of sorted arrays.o

Given k sorted arrays, we can use a Heap to perform a sorted traversal of all elements of all arrays.
We can push the first element of each sorted array to a minheap to get the overall minimum.
Important to note that we have to keep track of which array the element came from.
We can then remove the top element from the minheap and push the next element of the same array into the minheap.
This process can be repeated to get a sorted traversal of all elements.
