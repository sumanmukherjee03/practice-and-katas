### dfs

This is depth first traversal of a tree/graph.

Common approach is to use recursion or a stack to keep track of the parent node while traversing.
Space complexity will be O(H) where H is the max height of the tree/graph, ie at any point max there will
be H elements in the stack.

This traversal is generally useful in problems that require you to traverse from root to leaf.
