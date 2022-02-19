from collections import deque

def describe():
    desc = """
Problem : Topological Sort of a directed graph (a graph with unidirectional edges) is a linear ordering of its vertices such that
for every directed edge (U, V) from vertex U to vertex V, U comes before V in the ordering.
Given a directed graph, find the topological ordering of its vertices.

Example :
    Input: Vertices=4, Edges=[3, 2], [3, 0], [2, 0], [2, 1]
    Output: Following are the two valid topological sorts for the given graph:
        1) 3, 2, 0, 1
        2) 3, 2, 1, 0

    Input: Vertices=5, Edges=[4, 2], [4, 3], [2, 0], [2, 1], [3, 1]
    Output: Following are all valid topological sorts for the given graph:
        1) 4, 2, 3, 0, 1
        2) 4, 3, 2, 0, 1
        3) 4, 3, 2, 1, 0
        4) 4, 2, 3, 1, 0
        5) 4, 2, 0, 3, 1

-------------------------
    """
    print(desc)

#  If there is an edge from U to V then U<=V, ie, U comes before V in the ordering.
#  -  Source is any node that has no incoming edge and has only outgoing edges.
#  -  Sink is any node that only has incoming nodes and no outgoing edges.
#  -  A topological ordering starts with a source and ends at a sink.
#  -  A topological ordering is only possible when the graph has no directed cycles, ie the graph is a Direct Acyclic Graph (DAG)
#
#  Process to find topological sort
#  - Traverse graph in BFS
#  - Start with all the sources and in a stepwise fashio save all sources to a sorted list.
#  - Then remove all sources and their edges from the graph.
#  - This gives us new sources.
#  - Repeat the process above until all vertices are visited

#  Implementation
#    - Initialization
#      - Store the graph in an AdjacencyList, which means each parent vertex will have a list containing all of it's children.
#        This can be represented using a HashMap where the key will be a parent vertex number and the value will be a List containing child vertices
#      - To find the sources we will maintain another HashMap to count the in-degrees, ie, count the incoming edges of each vertex. Any vertex with 0 in-degree will be a source.
#    - Build the graph
#      - We will build the graph using the input and populate the AdjacencyList and build the in-degrees HashMap
#    - Find sources
#      - All vertices with in-degrees 0 will be the sources and we will be storing those in a Queue
#    - Sort
#      - For each source do the following
#        - Add it to the sorted List
#        - Get all it's children from the graph
#        - Decrement the in-degree of each child by 1
#        - If a childs in-degree becomes 0, add it to the Queue
#      - Repeat the previous step until the Queue becomes empty

#  Time complexity of this algo is O(V+E) where V is number of vertices and E is the number of edges.
#  Space complexity is also the same.
def topological_sort(numOfVertices, edges):
    sortedOrder = []
    if numOfVertices <= 0:
        return sortedOrder

    #  Remember that the edges have vertices starting with 0 index

    #  Initialize inDegrees HashMap to 0 for each vertex
    inDegrees = {v: 0 for v in range(0, numOfVertices)}

    #  Initialize inDegrees HashMap to 0 for each vertex
    adjacencyList = {v: [] for v in range(0, numOfVertices)}

    #  Now, go through the edges array and form the adjacency matrix of the graph
    for edge in edges:
        src, dest = edge
        adjacencyList[src].append(dest)
        inDegrees[dest] += 1

    q = deque()
    sources = [v for v in inDegrees if inDegrees[v] == 0] # Find all vertices that have in-degree 0 because they are the sources
    for s in sources:
        q.append(s)

    while q:
        # Get the number of element in the queue, because this represents the number of sources at each level
        numOfSources = len(q)
        for i in range(0, numOfSources):
            s = q.popleft() # pop a source from the front of the deque, like a queue
            sortedOrder.append(s)
            #  Iterate over the children of the source s
            for v in adjacencyList[s]:
                inDegrees[v] -= 1 # Reduce the in-degree of each child element by 1
                #  If the new in-degree of a child is 0, then that becomes a new source. So, append it to the queue
                if inDegrees[v] == 0:
                    q.append(v)

    #  IMPORTANT NOTE : If there is a cycle then the length of the sortedOrder array will be less than numOfVertices
    #  This condition can be used to determine if a directed graph has a cycle or not
    #  This is because where there is a cycle, we wouldnt be able to find any sources, so the queue would be empty.
    #  We would only be able to sort the vertices before the cycle.
    if len(sortedOrder) != numOfVertices:
        return []

    return sortedOrder

def main():
    describe()

    vertices = 4
    edges = [[3, 2], [3, 0], [2, 0], [2, 1]]
    print("Input : " + str(vertices) + ", " + str(edges))
    print("Topological sort: " + str(topological_sort(vertices, edges)))
    print("\n\n")

    vertices = 5
    edges = [[4, 2], [4, 3], [2, 0], [2, 1], [3, 1]]
    print("Input : " + str(vertices) + ", " + str(edges))
    print("Topological sort: " + str(topological_sort(vertices, edges)))
    print("\n\n")

    vertices = 7
    edges = [[6, 4], [6, 2], [5, 3], [5, 4], [3, 0], [3, 1], [3, 2], [4, 1]]
    print("Input : " + str(vertices) + ", " + str(edges))
    print("Topological sort: " + str(topological_sort(vertices, edges)))
    print("\n\n")

main()
