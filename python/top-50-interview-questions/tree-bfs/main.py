#  Problem : Given a tree print it's output in BFS order, ie one level at a time going from left to right

class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right

#  In a binary tree, to access a node we need access to it's parent.
#  So, when printing in BFS, we print the root element first, but we need to store it's children somewhere
#  so that we dont lose access to them later on and we can pick them out again and start exploring.
#  The best data structure to do this is a queue.
#  So, print root -> push neighbours to queue -> dequeue from queue and explore
#  Time complexity and space complexity is O(n)
def bfs01(root):
    if not root:
        return
    queue = []
    queue.push(root)
    while len(queue) > 0:
        node = queue.pop(0)
        print(node.data)
        if node.left:
            queue.append(node.left)
        if node.right:
            queue.append(node.right)
    return


# Same function as above but written in a recursive manner
def bfs02(root):
    if not root:
        return
    def bfsRecursive(queue):
        if len(queue) == 0:
            return
        node = queue.pop(0)
        print(node.data)
        if node.left:
            queue.append(node.left)
        if node.right:
            queue.append(node.right)
        return
    bfsRecursive([root])
