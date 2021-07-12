class Tree:
    def __init__(self, data, left = None, right = None):
        self.data = data
        self.left = left
        self.right = right

# Parameters:
#  root: Tree
# return type: List[List[int]]

def getValuesByLevel01(root):
    res = {}
    def bfs(node, level = 0):
        if node is None:
            return
        if res.get(level):
            res[level].append(node.data)
        else:
            res[level] = [node.data]
        bfs(node.left, level+1)
        bfs(node.right, level+1)

    bfs(root)

    out = []
    for k in res.keys():
        out.append(res[k])

    return out



#  This solution involves pushing the node as well the level into the queue when traversing via bfs
#  Time complexity is O(n)
def getValuesByLevel02(root):
    if root is None:
        return []

    queue = []
    out = []
    queue.append([root, 0])

    while len(queue) > 0:
        elm = queue.pop(0)
        node = elm[0]
        level = elm[1]

        if len(out) < level+1:
            out.append([node.data])
        else:
            out[level].append(node.data)

        if node.left:
            queue.append([node.left, level+1])
        if node.right:
            queue.append([node.right, level+1])

    return out
