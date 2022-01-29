class TreeNode():
    def __init__(self, val):
        self.value = val
        self.left = None
        self.right = None

def describe():
    desc = """
Problem : Given a binary tree where each node can only have a digit (0-9) value,
each root-to-leaf path will represent a number. Find the total sum of all the numbers represented by all paths.


---------------
    """
    print(desc)

def get_num_recurse(root, num, result):
    if root is None:
        return

    num = num * 10 + root.value
    if root.left is None and root.right is None:
        result[0] += num
    else:
        get_num_recurse(root.left, num, result)
        get_num_recurse(root.right, num, result)

    num = (num - root.value)/10


#  Time complexity is O(n)
def find_sum_of_path_numbers(root):
    #  Use an array here to hold the sum because python passes around arrays as reference and not value
    #  So it is easier to handle during recursion
    result = [0]
    get_num_recurse(root, 0, result)
    return result[0]

################ ALTERNATE IMPLEMENTATION #####################

#  At any point in the recursion, return the path's sum from that node onwards, ie the subtree from that node.
#  And that value is the sum of paths from left child subtree and right child subtree.
#  For an empty node that sum is 0
def find_root_to_leaf_path_numbers(currentNode, pathSum):
  if currentNode is None:
    return 0

  # calculate the path number of the current node
  pathSum = 10 * pathSum + currentNode.value

  # if the current node is a leaf, return the current path sum
  if currentNode.left is None and currentNode.right is None:
    return pathSum

  # sum of paths in a tree is sum of paths in it's left subtree plus it's right subtree
  # so, traverse the left and the right sub-tree and add the value returned from them
  return find_root_to_leaf_path_numbers(currentNode.left, pathSum) + find_root_to_leaf_path_numbers(currentNode.right, pathSum)

def alternate_find_sum_of_path_numbers(root):
    find_root_to_leaf_path_numbers(root, 0)


################### ENTRYPOINT #####################

def main():
    describe()
    root = TreeNode(1)
    root.left = TreeNode(0)
    root.right = TreeNode(1)
    root.left.left = TreeNode(1)
    root.right.left = TreeNode(6)
    root.right.right = TreeNode(5)
    print("Total Sum of Path Numbers: " + str(find_sum_of_path_numbers(root)))

main()
