### PROPERTIES OF A BST

1. All nodes on left subtree are smaller than the root, and all nodes on the right subtree are greater than the root
2. A child of a BST is also a BST
3. Inorder traversal of a BST produces a sorted array
4. We can insert, delete and find nodes in O(h) time, h being the height of the subtree


### PROPERTIES OF A PALINDROME

1. Reverse of a palindrome is also a palindrome
2. A palindrome plus the same char attached at the beginning and end of the palindrome is also a palindrome
3. In a palindrome every character appears an even number of times except for 1 char which can come odd number of times.

### Link Linked List

1. For reversing a linked list if there is only a next pointer
  - maintain 3 variables previous, current and next


### Double Link Linked List

1. For reversing a linked list if there are both next and previous pointers
  - maintain 2 pointers - one from beginning and another from end
  - swap values of nodes at position
  - move the pointers towards each other until they cross

### Permutations of a string

1. To count the number of possible permutations of a string you can use a mathematical formula
  - Given a string "cabcbcacc" -> {a: 2, b: 2, c: 5} -> 9!/(2! x 2! x 5!)

### BFS and DFS

_DFS_

In DFS the idea is go as deep as possible in a branch which pushing nodes onto a stack until you reach a leaf or if there's nothing left to explore any more.
Then you mark it as visited and backtrack by popping from the stack.

Stack is a good DS to keep track of that.
Recursion is also generally useful since recursion creates a stack

_BFS_

In BFS, the idea is to traverse level by level.
Queue is a good DS to keep track of nodes.

Visit a node first, then push it's children or neighbours into a queue.
Then dequeue from the queue and mark that visited and push it's neighbours into a queue.
Continue dequeuing until the queue is empty.
