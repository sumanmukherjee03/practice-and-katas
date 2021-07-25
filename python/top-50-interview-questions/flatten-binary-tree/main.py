#  Problem : Given a binary tree create a function that flattens it to a linked list IN PLACE
#  by following the preorder traversal.
#  For ex : [1,2,5,3,4,6]
#                     1
#                      /\
#                    /   \
#                  2      5
#                 /\      /\
#               /    \  /    \
#              3     4 6     null
#
#                TO
#
#              1
#              /\
#            /   \
#           null  2
#                 /\
#                /   \
#               null  3
#                     /\
#                    /   \
#                   null  4
#                         /\
#                        /   \
#                       null  5
#                             /\
#                            /   \
#                           null  6
#
#
# Output : [1, null, 2, null, 3, null, 4, null, 5, null, 6]

#  The easy solution for converting a tree to a linked list is traversing in preorder and pushing the data into a linked list
#  But that does not meet our constraint of flattening the tree in place.

class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right

#  The time complexity of this solution is O(n^2) because we are traversing the nodes once to flatten
#  and once to find the last node of the flattened subtree when in preorder
def flattenTree01(root):
    def isLeaf(node):
        if not node is None:
            if node.left is None and node.right is None:
                return True
            else:
                return False
        return True

    # Find the leaf node on any tree that would come last in the inorder traversal from that node.
    # Recursively keep going right if there is a right subtree
    # otherwise go left and then again try going right
    # If there are no right subtrees at all, then it'll be the last leaf of the left subtree.
    def rightMostLeaf(node):
        if isLeaf(node):
            return node
        elif not node is None and not node.right is None:
            return rightMostLeaf(node.right)
        elif not node is None and not node.left is None:
            return rightMostLeaf(node.left)
        else:
            return node

    # Keep the current right subtree in a temp var
    # Flatten the left subtree and make it the right subtree of the node
    # Get the leaf node from the left subtree that comes last in the preorder traversal
    # Attach the temp var containing the original right subtree as the right subtree of the leaf node from above
    # Flatten the original right subtree you just attached in the step above
    # Recursively perform this operation
    def flattenRight(node):
        if node is None:
            return
        temp = node.right
        flattenRight(node.left)
        rNode = rightMostLeaf(node.left)
        node.right = node.left
        node.left = None
        if not rNode is None:
            rNode.right = temp
        else:
            node.right = temp
        flattenRight(temp)

    flattenRight(root)




#  A much simpler solution - same solution as above but a much simpler code
#  Time complexity is O(n^2)
def flattenTree02(root):
    if root is None:
        return
    else:
        # Flatten the left subtree
        # Flatten the right subtree
        # Make the left subtree as the right subtree
        # Make the left subtree as null
        # Traverse to the end of the right subtree to find the last node
        # Attach the original flattened right subtree as the right child of the node found in the previous step
        flattenTree02(root.left)
        flattenTree02(root.right)
        rSubtree = root.right
        root.right = root.left
        root.left = None
        n = root
        while n is not None:
            n = n.right
        n.right = rSubtree



#  There is also a O(n) solution but i found it difficult to understand
