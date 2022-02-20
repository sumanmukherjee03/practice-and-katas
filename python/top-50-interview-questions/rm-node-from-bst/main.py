#  Problem : Given a binary search tree root and a num, create a function that deletes the node that contains the num then returns the root.

#  ----------------------------

class Tree():
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right


def inorderTraversal(root):
    return (inorderTraversal(root.left) + [root.data] + inorderTraversal(root.right)) if root else []


def getMinNode(node):
    while node.left is not None:
        node = node.left
    return node

#  NOTE : When traversing inorder all the elements of a BST are sorted.
#  All elements in left subtree are smaller than root and all element in right subtree are larger than root.
#  So, when deleting a node from a BST, we must maintain the sorted order from the inorder traversal
#  ie, the node will be replaced either by the element right before that node in inorder traversal or the element right after it in inorder traversal
#  Based on that knowledge, we can simply find the next smallest node and replace the node we found with the data/value of the next smallest node
#  Then recursively delete the next smallest node from the right subtree of the found node
#
#  Time complexity is O(h)
def deleteNodeBst(root, num):
    if root is None:
        return None
    elif num < root.data:
        root.left = deleteNodeBst(root.left)
    elif num > root.data:
        root.right = deleteNodeBst(root.right)
    else:
        #  When the node we are looking to delete is found, there are 3 cases to consider
        #  If the node to be deleted has no left subtree then simply remove that node and return it's right subtree
        #  If the node to be deleted has no right subtree then simply remove that node and return it's left subtree
        #  Otherwise if there is both a left and a right subtree, then we want to find the next minimum node in the in-order traversal
        #  and replace the data of current node with the successor node's data.
        #  Then from the right subtree delete the successor node recursively. This will readjust the right subtree to become a BST again.
        if root.left is None:
            return root.right
        elif root.right is None:
            return root.left
        else:
            successor = getMinNode(root.right)
            root.data = successor.data
            root.right = deleteNodeBst(root.right, successor.data)
    return root
