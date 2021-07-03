# Time complexity of DFS in tree traversal is O(n) where n is the number of nodes and space complexity is O(h)
#  where h is the height of the tree

class Tree:
  def __init__(self, data, left = None, right = None):
    self.data = data
    self.left = left
    self.right = right

def dfsPreorder(root):
    res = []
    if root is None:
        return res
    res.append(root.data)
    if root.left:
        res = res + dfsPreorder(root.left)
    if root.right:
        res = res + dfsPreorder(root.right)
    return res

def dfsInorder(root):
    res = []
    if root is None:
        return res
    if root.left:
        res = res + dfsPreorder(root.left)
    res.append(root.data)
    if root.right:
        res = res + dfsPreorder(root.right)
    return res

def dfsPostorder(root):
    res = []
    if root is None:
        return res
    if root.left:
        res = res + dfsPreorder(root.left)
    if root.right:
        res = res + dfsPreorder(root.right)
    res.append(root.data)
    return res
