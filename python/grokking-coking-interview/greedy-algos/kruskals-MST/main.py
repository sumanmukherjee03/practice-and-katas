from collections import defaultdict

def describe():
    desc = """
Minimum Spanning Tree description :
-----------------------------------
Given a connected and undirected graph, a spanning tree of that graph is a subgraph
that is a tree and connects all the vertices together.
A single graph can have many different spanning trees.
A minimum spanning tree (MST) or minimum weight spanning tree for a weighted, connected, undirected graph
is a spanning tree with a weight less than or equal to the weight of every other spanning tree.
The weight of a spanning tree is the sum of weights given to each edge of the spanning tree.
A minimum spanning tree has V-1 edges where V is the number of vertices in the given graph.

Problem : Find the minimum spanning tree of a graph using Kruskals algo

---------------------
    """
    print(desc)

#  Steps to find MST using Kruskals algo :
#  1. Sort all edges in non-decreasing order according to their weight
#  2. Pick the smallest edge and check if it forms a cycle with the spanning tree so far
#  .   - If there is no cycle then include this edge
#  3. Repeat step 2 until there are (V-1) edges in the spanning tree

class Graph:
    def __init__(self, vertices):
        self.V = vertices # V represents the number of vertices in the graph
        #  defaultdict takes a function to invoke for returning a default value when a key is not found
        self.graph = [] # The graph is represented as a dictionary with each element being like a tuple of <node 1, node 2, weight of edge>

    def addEdge(self, u, v, w):
        self.graph.append([u, v, w])



    #  Each index of the parents array represents a vertex/node.
    #  The value of the parents array at each index represents which set the vertex belongs to.
    #  ie, if the value at index 1 is 2, then that means that vertex 1 belongs to the same set as vertex 2.
    #  ie, parent of index 1 is index 2. So, now recursively go looking for parent of index 2.
    #  Since the parents array is initialized with -1, if we hit the case that the parent at index x is -1
    #  it means there are no more parents to explore and this is the end of the chain.
    def find_parent(self, parent, i):
        if parent[i] == -1:
            return i
        if parent[i] != -1:
            return self.find_parent(parent, parent[i])

    #  If 2 vertices are in disjoint sets and there exists a vertex connecting the 2, then union them
    #  ie, join vertex x to the spanning tree of y or join vertex x to the same set as y
    def union(self, parent, x, y):
        parent[x] = y

    def kruskalMST(self):
        result = []
        i = 0 # A pointer for traversing the sorted edged
        e = 0 # A pointer used for result
        # get a copy of the graph which is sorted based on weights of the edges
        graph = sorted(self.graph, key = lambda x : x[2])
        parent = [-1] * (self.V)

        #  Keep iterating and taking edges into the result set as long as there are not V-1 edges
        while e < (self.V - 1):
            u, v, w = graph[i]
            i += 1
            x = self.find_parent(parent, u)
            y = self.find_parent(parent, v)
            if x != y:
                e += 1
                result.append([u, v, w])
                self.union(parent, x, y)

        return result


def main():
    describe()

    g = Graph(4)
    g.addEdge(0, 1, 10)
    g.addEdge(0, 2, 6)
    g.addEdge(0, 3, 5)
    g.addEdge(1, 3, 15)
    g.addEdge(2, 3, 4)

    res = g.kruskalMST()
    print(str(res))

main()
