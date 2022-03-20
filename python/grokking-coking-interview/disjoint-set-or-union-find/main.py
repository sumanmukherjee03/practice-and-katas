from collections import defaultdict

def describe():
    desc = """
Disjoint-Set and Union-find :
-----------------------------
A disjoint-set data structure is a data structure that keeps track of a set of elements partitioned into a number of disjoint (non-overlapping) subsets.
A union-find algorithm is an algorithm that performs two useful operations on such a data structure:
    - Find: Determine which subset a particular element is in. This can be used for determining if two elements are in the same subset.
    - Union: Join two subsets into a single subset. Here first we have to check if the two subsets belong to same set. If no, then we cannot perform union.


Prolem : Check whether a given graph contains a cycle or not

    """
    print(desc)


class Graph:
    def __init__(self, vertices):
        self.V = vertices # no of vertices
        #  defaultdict takes a function to invoke for returning a default value when a key is not found
        self.graph = defaultdict(list) # The graph is represented as a dictionary

    #  Think of adding an edge in the graph as going from v1 -> [v2].
    #  If there are multiple in the list, as in, v1 -> [v2, v3]
    #  Then it means we have 2 edges v1 -> v2, v1 -> v3
    #  We have one egde for any two vertex i.e 1-2 is either 1 -> 2 or 2 -> 1 but not both.
    #  self.graph is what represents the adjacency matrix
    def addEdge(self, u, v):
        self.graph[u].append(v)

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

    #  Initially all elements of the parents array is -1.
    #  Each index of the parents array represents a vertex.
    #  Initially consider each vertex is in a separate set.
    #  As we process the edges by traversing the adjacency matrix, update the parents.
    #  If we, process an edge (0,1), then parent becomes 1,-1,-1.
    #  If we process an edge now with (1,2), then parent becomes 1, 2, -1. ie 0th vextex is connected to 1 and 1 is connected to 2 and so on.
    #  So, if we take 2 arbitrary vertices and they traverse up their parents recursively and end up in the same node,
    #  that means they belong to the same set and hence we have a cycle.
    #  If the parents do not end up as the same one, then we union the nodes, as in consider the edge in the acyclic spanning tree that we are built so far.
    def isCyclic(self):
        parent = [-1] * (self.V)
        #  For each vertex, traverse edges it is connected to and keep forming the spanning tree.
        #  If the 2 vertices we are considering in this edge belong to disjoint sets, ie they have different parents then union them
        #  NOTE : UNION THE PARENTS, NOT THE NODES THEMSELVES
        #  Otherwise there is a cycle
        for i in self.graph:
            for j in self.graph[i]:
                x = self.find_parent(parent, i)
                y = self.find_parent(parent, j)
                if x == y:
                    return True
                self.union(parent,x,y)



def main():
    describe()

    g = Graph(3)
    g.addEdge(0, 1)
    g.addEdge(1, 2)
    g.addEdge(2, 0)
    print(str(g.graph))
    res = g.isCyclic()
    if res is True:
        print("Graph is cyclic")
    else:
        print("Graph is acyclic")
    print("\n\n")

    g = Graph(4)
    g.addEdge(0, 1)
    g.addEdge(0, 2)
    g.addEdge(2, 3)
    print(str(g.graph))
    res = g.isCyclic()
    if res is True:
        print("Graph is cyclic")
    else:
        print("Graph is acyclic")
    print("\n\n")
main()
