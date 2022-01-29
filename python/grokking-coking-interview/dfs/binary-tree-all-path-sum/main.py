class TreeNode():
    def __init__(self, val):
        self.value = val
        self.left = None
        self.right = None

def describe():
    desc = """
Problem : Given a binary tree and a number S, find all paths from root-to-leaf such that the sum of all the node values of each path equals S

--------------
    """
    print(desc)


#  currentPath maintains the current state of the pointer in the callstack, as in, how deep we are down a path
#  allPaths is being passed around like a global
#  Time complexity of this is O(n^2)
def find_paths_recursive(root, target, currentPath, allPaths):
    #  This guard clause handles the case when the left child or right child is None for a node and this func has been called
    if root is None:
        return

    currentPath.append(root.value)
    if root.left is None and root.right is None and root.value == target:
        #  If a path is found at this point append it to the list of possible paths
        #  But do not return from here. Continue on in the function. You want to remove this node from
        #  the currrentPath before you return so that other paths can also be explored.
        #  Remember that recursion does not stop as soon as a match is found.
        allPaths.append(list(currentPath)) # This is just a nuisance i think - having to convert the currentPath to a list before appending
    else:
        find_paths_recursive(root.left, target - root.value, currentPath, allPaths)
        find_paths_recursive(root.right, target - root.value, currentPath, allPaths)

    #  If the path led to a successful search, then that currentPath would have been appended to allPaths by now.
    #  This statement is simply a reverse of currentPath.append(root.value) before this function returns in the recursive callstack.
    #  Say we found a path [12, 7, 4] and we are in the call stack for the node with value 4. Then that path
    #  would have been added to allPaths. Now, we can remove 4 from current path and return the call stack,
    #  ie currentPath would be [12, 7]. And the call returns to the parent node of 4, so that it can continue down the path of the other child.
    del currentPath[-1]

def find_paths(root, s):
    allPaths = []
    find_paths_recursive(root, s, [], allPaths)
    return allPaths

def main():
    describe()
    root = TreeNode(12)
    root.left = TreeNode(7)
    root.right = TreeNode(1)
    root.left.left = TreeNode(4)
    root.right.left = TreeNode(10)
    root.right.right = TreeNode(5)
    s = 23
    print("Tree paths with sum " + str(s) +
            ": " + str(find_paths(root, s)))

main()
