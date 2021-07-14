#  Problem : Given an undirected graph of integers represented by an adjacency list and
#  an integer root, create a function that prints its values in DFS, by considering whose
#  value is root as the arbitrary node.
#  For ex : {"5": [8,1,2], "8": [5,12,14,4], "12": [5,8,14], "14": [8,12,4], "4": [8,14], "1": [5,7], "7": [1,16], "16": [7]}, root = 5
#  Output : 5 8 1 12 14 4 7 16

#  Remember that there are 2 ways of representing graphs, using an adjacency list or adjacency matrix.
#  This representation shown above is using an adjacency list.
#  Using an adjacency matrix the number of rows and columns is the same as the number of nodes, so n x n matrix.
#  And each cell (i,j) of a matrix represents if there is an edge from i -> j.
#  In an undirected graph if there is a 1 in (i,j) there will also be a 1 in (j,i), ie the matrix will be symetrical in nature.
#  Similarly for the adjacency list representation, you will notice that if j is a neighbour of i, i is also listed as a neighbour of j.

#  In adjacency list representation of an undirected graph, the number of elements in adjacency lists is 2 * (number of vertices)

#  Time complexity of this traversal is O(V+E)

class Graph:
    def __init__(self, adjList = {}):
        # the adjacency list is of type Dict[int,List[int]]
        self.adjList = adjList

#  The time complexity is O(V+E)
def bfs(graph, root):
    traversal = []
    queue = [root]
    visited = {root}

    # Continue while loop as long as there are elements in the queue
    while len(queue) > 0:
        node = queue.pop(0) # Get the first element from the queue, ie dequeue
        traversal.append(node) # This node is traversed, so push it to the traversal array for printing later
        # Iterate over the neighbours and if a node is not already visited, push it into the queue for visiting later
        for n in graph.adjList[str(node)]:
            if n not in visited:
                queue.append(n)
                #  Mark the neighbour as visited, by adding the element to the set
                #  We mark it as visited here and not while printing it because a node can
                #  have many incoming vertices. So, it will be listed as a neighbour for multiple other nodes before it gets printed.
                #  but we only want to push it to the queue once. That's why it needs to be marked as visited here.
                visited.add(n)

    print(traversal)
    return
