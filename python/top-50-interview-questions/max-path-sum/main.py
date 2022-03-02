def describe():
    desc = """
Problem : Given a non-empty binary tree root, return the maximum path sum.
          Note that for this problem, a path goes from one node to another by traversing edges.
          The path must have at least one edge and it does not have to pass by the root.

NOTE : When we say that the path does not have to pass through the root, it implies that a subtree can also have the max path.

-------------------
    """
    print(desc)

class Tree:
    def __init__(self, data, left=None, right=None):
        self.data = data
        self.left = left
        self.right = right

def maxPathSum(root):
    globalMaxSum = [float("-inf")]
    dfs(root, globalMaxSum)
    return globalMaxSum[0]

#  Time complexity is O(n)
def dfs(root, globalMaxSum):
    if root is None:
        return float("-inf")
    else:
        #  We need to find the max paths when we go to the left and when we go to the right
        #  So, we store those 2 values in variables left and right
        left = dfs(root.left, globalMaxSum)
        right = dfs(root.right, globalMaxSum)
        #  Now, we have 2 types of max sum. One, where the path originates from the top and goes down a subtree, ie left or right.
        #  And another possibility is one where the path does not originate from the root.
        #  In that case the path can go through both left and right subtrees.

        #  When we consider the situation where the path is from root and goes down a subtree, there are 3 possibilities only
        #  - data at root is a higher value than left or right subtree because there can be negative nodes
        #  - data at root + left subtree is high value
        #  - data at root + right subtree is high value
        maxFromTop = max(root.data, root.data + left, root.data + right)

        #  When we consider the situation where the path does not originate from root and can go from left subtree to the right subtree, there are 4 possibilities
        #  - data at root is a higher value than left or right subtree because there can be negative nodes
        #  - data at root + left subtree is high value
        #  - data at root + right subtree is high value
        #  - data at root + max from left subtree + max from right subtree is a higher value
        maxNoTop = max(maxFromTop, root.data + left + right)

        #  And finally we keep track of the global sum as we recurse
        #  It is the maximum between whatever is the current global max sum and the max sum of this current subtree
        #  Also, note that we use an array here to pass global max sum around by reference
        globalMaxSum[0] = max(globalMaxSum[0], maxNoTop)

        #  Finally return maxFromTop because that's what we need when we backtrack to the parent at any stage.
        #  This value is used above by the recusive calls to get the max from top for left and right subtrees.
        return maxFromTop

def main():
    describe()

    root = Tree(11)
    root.left = Tree(1)
    root.right = Tree(2)
    root.left.left = Tree(4)
    root.left.left.right = Tree(2)
    root.right.left = Tree(5)
    root.right.right = Tree(10)
    root.right.left.right = Tree(8)
    print("Result : " + str(maxPathSum(root)))

main()
