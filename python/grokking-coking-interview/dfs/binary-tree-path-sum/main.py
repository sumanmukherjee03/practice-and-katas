from collections import deque

class TreeNode():
    def __init__(self, val):
        self.value = val
        self.left = None
        self.right = None

def describe():
    desc = """
Problem : Given a binary tree and a number S, find if the tree has a path from root-to-leaf
          such that the sum of all the node values of that path equals S.

---------------

    """
    print(desc)

#  Time complexity is O(n) where n is the number of nodes in the tree
def has_path(root, target):
    #  If the root is empty return False
    #  If at any point the value of root is  larger than the current target then also return false
    if root is None or root.value > target:
        return False

    #  If the current node is a leaf node and is of the same value as the target then we have found a path
    if root.left is None and root.right is None and root.value == target:
        return True

    #  Otherwise recursively search for the remainder in the left or right subtree
    target = target - root.value
    return has_path(root.left, target) or has_path(root.right, target)


#  Think of this traversal as - at each node it puts nodes from right to left into the stack. That
#  way the leftmost nodes are at the top of the stack. Then it picks up and process nodes from the left, ie the top of the stack.
#  This process essentially also ensure that at any level the nodes in the stack are from right to left in the stack, ie the left most in that level will get popped and processed first
def dfs_traversal(root):
    visited = []
    stack = deque()
    stack.append(root)
    while stack:
        node = stack.pop() # Always be popping the top of the stack
        visited.append(node.value)
        # For a tree DFS traversal, push the right node first so that it gets popped from the stack first
        #  At any point the stack has newer nodes at the right hand end and those will get popped first and that is desirable
        if node.right is not None:
            stack.append(node.right)
        if node.left is not None:
            stack.append(node.left)
        #  print(visited, map(lambda x: x.value, stack))
    return visited



def main():
    describe()
    root = TreeNode(12)
    root.left = TreeNode(7)
    root.right = TreeNode(1)
    root.left.left = TreeNode(9)
    root.right.left = TreeNode(10)
    root.right.right = TreeNode(5)
    print("Tree has path: " + str(has_path(root, 23)))
    print("Tree has path: " + str(has_path(root, 16)))
    print("Tree has DFS: ", dfs_traversal(root))

main()
