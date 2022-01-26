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

main()
