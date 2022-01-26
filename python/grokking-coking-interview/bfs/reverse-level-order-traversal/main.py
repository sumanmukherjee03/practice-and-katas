from collections import deque

class TreeNode():
    def __init__(self, val):
        self.value = val
        self.left = None
        self.right = None

def describe():
    desc = """
Problem : Given a binary tree, populate an array to represent its level-by-level traversal in reverse order
          i.e., the lowest level comes first. You should populate the values of all nodes in each level from left to right in separate sub-arrays.

--------------
    """
    print(desc)

#  This is almost the same as the normal level by level bfs
#  except that the output of each level is appended to the left of resulting queue
#  Time complexity is also of O(n)
def traverse(root):
    result = deque()
    if root is None:
        return result
    queue = deque()
    queue.append(root)
    while queue:
        level = []
        #  Only visit the length of the queue number of elements because that represents the number of nodes in each level
        for _ in range(0,len(queue)):
            #  Take out from the front of the queue, ie the left and insert at the end of the queue, ie right
            node = queue.popleft()
            if node.left is not None:
                queue.append(node.left)
            if node.right is not None:
                queue.append(node.right)
            level.append(node.value)
        result.appendleft(level)
    return result

def main():
    describe()
    root = TreeNode(12)
    root.left = TreeNode(7)
    root.right = TreeNode(1)
    root.left.left = TreeNode(9)
    root.right.left = TreeNode(10)
    root.right.right = TreeNode(5)
    print("Reverse level order traversal: " + str(traverse(root)))

main()


