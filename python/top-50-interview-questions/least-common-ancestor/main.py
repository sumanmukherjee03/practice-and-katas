#  Problem : Given a binary tree root and 2 integers num1 and num2, create a function
#  that returns the lowest common ancestor of the 2 integers in the tree. The LCA is
#  the deepest root node that has both num1 and num2 as descendents and we consider a node
#  as a descendent of itself. All values are unique and assume that num1 and num2 always exist in the tree.

class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right

#  Time complexity is O(n)
def lowestCommonAncestor01(root, num1, num2):
    #  Get the path from root to the node with the num as value and return if a path was found or not
    #  Else return false
    #  We maintain a stack to push nodes into so that we can generate the path to find a node.
    #  If the node is found via going through that node, keep it in the stack, otherwise pop
    #  from the stack after exploring the left and right subtrees
    def getPath(root, path, num):
        if root is None:
            return False
        #  Append root to path because we are going to be exploring this as a potential node in the path
        path.append(root)
        # If found the node we have the path
        if root.data == num:
            return True
        # If found the node in left or right subtree return true
        if getPath(root.left, path, num) or getPath(root.right, path, num):
            return True
        #  If the node wasnt found in the left or right subtree, it means this node is not in the path
        path.pop()
        return False

    path1 = []
    path2 = []
    if not getPath(root, path1, num1) or not getPath(root, path2, num2):
        return None

    # Find the length of the minimum path and keep traversing the paths as long as the nodes match
    # As soon as the nodes mismatch return the last matched node because that's the parent.
    minLen = min(len(path1), len(path2))
    i = 0
    while i < minLen:
        if path1[i].data == path2[i].data:
            i += 1
        else:
            break
    return path1[i-1]



#  There is another solution to this problem but that is harder to visualize
