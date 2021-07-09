#  Problem : Given a binary tree root check that it is a valid binary search tree
#  A BST is valid when every node must be strictly greater than all nodes on the left
#  and strictly smaller than all nodes on the right
#  Or in other words the root must be bigger than all nodes in the left subtree and all nodes in the right subtree
#  Ex : Input - [16, 8, 22, 3, 11, null, null, 1, 6] -> this represents a level order traversal (BFS) array
#  Output : True

class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right

# Every node in a BST has value in the range of (it's closest left parent, it's closest right parent)
# For a node that is on the right subtree will have it's own parent as it's closest left parent, so that will be the min.
# And it's closest right parent will be whatever is the closest right parent of it's own parent.
# Similarly, the closest right parent of a node which is a left child will be it's own parent.
# And it's closest left parent will be whatever is the closest left parent of it's own parent.
# In this specific problem we do not need the parent nodes per se, ie it is not a problem where the links are of importance.
# In this problem the values from the left and right parents are important, not the information which node it came from.
#
# Using that property above the time complexity of this is O(n)
def isBst01(root, min = float("-inf"), max = float("inf")):
    if root is None:
        return True
    elif root.data <= min or root.data >= max:
        return False
    else:
        return isBst01(root.left, min, root.data) and isBst01(root.right, root.data, max)


#  Another property of BST that we can use is that the inorder traversal of a BST is a sorted array
def isBst02(root):
    if root is None:
        return True

    def inorder(node):
        if node is None:
            return []
        return inorder(node.left) + [node.data] + inorder(node.right)

    # Here we search for duplicate elements as well because a BST shouldnt contain any duplicates
    def isSorted(arr):
        for i in range(len(arr)-1):
            if arr[i+1] <= arr[i]:
                return False
        return True

    inorderTraversal = inorder(root, [])
    print(inorderTraversal)
    return isSorted(inorderTraversal)


#  Same logic as the one before but in a slightly more terse implementation.
#  If left node is a BST with it's precedent being the same as the current node.
#  Check the current node is bigger than it's predecent.
#  Check the right node is a BST with it's precedent being the current node.
def isBst03(root):
    def isBstRecursive(node, precedent):
        if node is None:
            return True
        if not isBstRecursive(node.left, precedent):
            return False
        if node.data <= precedent[0]:
            return False
        else:
            precedent[0] = node.data
        if not isBstRecursive(node.right, precedent):
            return False
        return True

    return isBstRecursive(root, [float("-inf")])
