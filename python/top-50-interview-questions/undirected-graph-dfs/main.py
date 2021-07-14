#  Problem : Given an undirected graph of integers represented by an adjacency list and
#  an integer root, create a function that prints its values in DFS, by considering whose
#  value is root as the arbitrary node.
#  For ex : {"5": [8,1,2], "8": [5,12,14,4], "12": [5,8,14], "14": [8,12,4], "4": [8,14], "1": [5,7], "7": [1,16], "16": [7]}, root = 5
#  Output : 5 8 12 14 4 1 7 16

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

def dfs01(graph, root):
    visited = {}
    traversal = []

    def dfs_recursive(g, r):
        visited[str(r)] = True
        traversal.append(r)
        neighbours = graph.adjList[str(r)]
        for n in neighbours:
            if not visited.get(str(n)):
                dfs_recursive(g, n)

    print(traversal)
    return


def dfs02(graph, root, visited=set()):
    if root in visited:
        return
    print(root)
    visited.add(root)
    for node in graph.adjList[str(root)]:
        dfs02(graph,node,visited)
    return
