#  Problem : Given a binary tree find out if it is balanced.
#  A balanced binary tree is one in which the difference between the height of left vs right subtree is at most 1.

class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right

#  Recursively calculate the diff of the height of the left and right subtree to find out if a node is balanced
#  The entire tree is considered balanced when all nodes appear balanced
#  The time complexity of this O(nlogn)
def isBalanced01(root):
    if root is None:
        return True

    #  Recursively calculate the height of a tree
    #  The time complexity of this function is O(n)
    def treeHeight(node):
        if node is None:
            return -1
        else:
            return 1 + max(treeHeight(node.left), treeHeight(node.right))

    return abs(treeHeight(root.left) - treeHeight(root.right)) <= 1 and isBalanced01(root.left) and isBalanced01(root.right)

#  In the previous solution we are calculating the height of each node and also calculating if each node is balanced or not.
#  This can be improved by calculating the balancing factor while calculating the height
#  Time complexity of this solution is O(n)
def isBalanced02(root):
    if root is None:
        return True

    #  In python to use a variable that's defined outside the function use a list
    #  because a list will be passed by reference and not by value
    res = [True]

    #  Recursively calculate the height of a tree
    #  But this time while calculating the height also check if the tree is balanced.
    #  The time complexity of this function is O(n)
    def _is_balanced(node, out):
        if root is None:
            return -1
        else:
            leftHeight = _is_balanced(node.left)
            rightHeight = _is_balanced(node.right)
            if abs(leftHeight - rightHeight) > 1:
                out[0] = False
            return 1 + max(leftHeight, rightHeight)

    _is_balanced(root, res)
    return res[0]
