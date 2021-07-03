class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right


#  The time complexity is O(n) where n is the number of nodes
#  and the space complexity is O(h) where h is the height of the tree
#  because that is what the recursion call stack will look like at it's max
def reverseTree(root):
    if root is None:
        return
    leftTree = root.left
    rightTree = root.right
    reverseTree(leftTree)
    root.right = leftTree
    reverseTree(rightTree)
    root.left = rightTree
