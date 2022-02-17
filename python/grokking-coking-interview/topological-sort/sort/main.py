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
#  - Repeat the process above.

def topological_sort(vertices, edges):
    sortedOrder = []
    # TODO: Write your code here
    return sortedOrder

def main():
    describe()
    print("Topological sort: " + str(topological_sort(4, [[3, 2], [3, 0], [2, 0], [2, 1]])))
    print("Topological sort: " + str(topological_sort(5, [[4, 2], [4, 3], [2, 0], [2, 1], [3, 1]])))
    print("Topological sort: " + str(topological_sort(7, [[6, 4], [6, 2], [5, 3], [5, 4], [3, 0], [3, 1], [3, 2], [4, 1]])))

main()
