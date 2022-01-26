from collections import deque

class TreeNode():
    """Docstring for TreeNode. """
    def __init__(self, val):
        self.value = val
        self.left, self.right = None, None


def describe():
    desc = """
Problem : Given a binary tree, populate an array to represent its level-by-level traversal.
          You should populate the values of all nodes of each level from left to right in separate sub-arrays.
--------------
    """
    print(desc)

#  Time complexity is O(n)
def traverse(root):
    result = []
    if root is None:
        return result

    queue = deque([root])
    while queue:
        # since you are increasing the size of the queue inside the nexted loop as you visit nodes
        #  the current length/size of the queue represents how many nodes there are at this current level
        levelsize = len(queue)
        level = []
        # perform visit operation by popping only the number of nodes from the queue that are on the current level
        for _ in range(0, levelsize):
            node = queue.popleft() # pop from the front of the queue
            level.append(node.value)

            # add left and right nodes to the end of the queue
            if node.left is not None:
                queue.append(node.left)
            if node.right is not None:
                queue.append(node.right)

        result.append(level)
    return result

def main():
    describe()
    root = TreeNode(12)
    root.left = TreeNode(7)
    root.right = TreeNode(1)
    root.left.left = TreeNode(9)
    root.right.left = TreeNode(10)
    root.right.right = TreeNode(5)
    print("Level order traversal: " + str(traverse(root)))

main()
