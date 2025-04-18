### PROPERTIES OF A BST

1. All nodes on left subtree are smaller than the root, and all nodes on the right subtree are greater than the root
2. A child of a BST is also a BST
3. Inorder traversal of a BST produces a sorted array and can be used to validate if a BST is valid or not
4. We can insert, delete and find nodes in O(h) time, h being the height of the subtree
5. Value of any node of a BST lies between the value of it's closest left parent and closest right parent


### PROPERTIES OF A PALINDROME

1. Reverse of a palindrome is also a palindrome
2. A palindrome plus the same char attached at the beginning and end of the palindrome is also a palindrome
3. In a palindrome every character appears an even number of times except for 1 char which can come odd number of times.

### Link Linked List

1. For reversing a linked list if there is only a next pointer
  - maintain 3 variables previous, current and next
  - in this case we need to modify links and not values, like in the case of a double link linked list
2. When detecting a cycle in a single link linked list with Floyd Cycle Detection algo
  - Remember that the node forming the cycle has 2 links which point to it as next
3. To find the middle of a linked list, use a slow and fast pointer. Move the fast pointer 2X the speed.
      That way by the time the fast pointer reaches the end the slow pointer has reached the middle only.
4. For palindrome, one solution with lower time complexity is to reverse the right half from the middle.
      Then compare element by element between the left half and the reversed right half.


### Double Link Linked List

1. For reversing a linked list if there are both next and previous pointers
  - maintain 2 pointers - one from beginning and another from end
  - swap values of nodes at position
  - move the pointers towards each other until they cross

### Permutations of a string

1. To count the number of possible permutations of a string you can use a mathematical formula
  - Given a string "cabcbcacc" -> {a: 2, b: 2, c: 5} -> 9!/(2! x 2! x 5!)


### String traversal

1. You can use a hash table to store the information if a char has been visited or not
2. You can maintain an array of size 128, each position representing an ascii code to store the index of last occurance of a char.
      This is usually helpful in sliding window technique.
3. For substring problems consider a dynamic sliding window technique with 2 pointers.
    - The stop pointer can move as an iteration
    - The start pointer can change position based on some condition


### Max subarray

1. Dynamic programming with Kadane's algo

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
