class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right


#  The time complexity is O(n) where n is the number of nodes
#  and the space complexity is O(h) where h is the height of the tree
#  because that is what the recursion call stack will look like at it's max
#  This implementation is based on DFS
def reverseTree01(root):
    if root is None:
        return
    leftTree = root.left
    rightTree = root.right
    reverseTree01(leftTree)
    root.right = leftTree
    reverseTree01(rightTree)
    root.left = rightTree



#  This implementation of reversing a binary tree is based on BFS
def reverseTree02(root):
    if root is None:
        return
    queue = [root]
    while len(queue) > 0:
        node = queue.pop(0)
        node.left, node.right = node.right, node.left
        if node.left:
            queue.append(node.left)
        if node.right:
            queue.append(node.right)
