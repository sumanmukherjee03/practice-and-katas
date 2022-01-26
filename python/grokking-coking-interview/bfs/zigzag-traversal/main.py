from collections import deque

class TreeNode():
    def __init__(self, val):
        self.value = val
        self.left = None
        self.right = None

def describe():
    desc = """
Problem : Given a binary tree, populate an array to represent its zigzag level order traversal.
          You should populate the values of all nodes of the first level from left to right,
          then right to left for the next level and keep alternating in the same manner for the following levels.

------------------
    """
    print(desc)

#  Time complexity is O(n)
def traverse(root):
    res = []
    if root is None:
        return res

    #  Dont change how you are performing the traversal, consider how you are building the result set
    level = 1
    queue = deque()
    queue.append(root)
    while queue:
        num_nodes_in_level = len(queue)
        nodes_at_level = deque() # Maintain a deque for containing the nodes of each level
        for _ in range(num_nodes_in_level):
            node = queue.popleft() # Always pop nodes from the front of the queue and append children to the end of the queue

            #  When inserting the popped node during the level traversal into the local container make the decision
            #  to insert at the front of the local container or end of the local container
            if level % 2 == 1:
                nodes_at_level.append(node.value)
            else:
                nodes_at_level.appendleft(node.value)

            if node.left is not None:
                queue.append(node.left)
            if node.right is not None:
                queue.append(node.right)

        level += 1
        res.append(list(nodes_at_level))
    return res

def main():
    describe()
    root = TreeNode(12)
    root.left = TreeNode(7)
    root.right = TreeNode(1)
    root.left.left = TreeNode(9)
    root.right.left = TreeNode(10)
    root.right.right = TreeNode(5)
    root.right.left.left = TreeNode(20)
    root.right.left.right = TreeNode(17)
    print("Zigzag traversal: " + str(traverse(root)))

main()
