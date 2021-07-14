#  Problem : Given a non-empty array of non-negative integers, where the value at each cell
#  represents the maximum jump that we can perform from that index, write a function to determine
#  if wer can reach the last index starting from the first one.
#  For ex : [3,0,2,0,2,1,4,3]
#  Output : true

#  The greedy method of jumping to the max is not always going to work because you can
#  end up in an index with value 0 and return false, whereas there might in fact be a path.

#  This is a classic backtracking problem and the time complexity here is O(2^n)
def canJump01(arr, i = 0):
    #  If we reached the end, great
    if i == len(arr)-1:
        return True
    #  Check all possibilities of a jump from an index based on the value in a cell
    for j in range(1,arr[i]+1):
        if canJump01(arr,i+j):
            return True
    return False


#  The solution above has a lot of repeated repeated calls when you look at the tree of recursive calls.
#  This increases the time complexity signnificantly and can be reduced with dynamic programming.
#  Time complexity is O(n^2)
def canJump02(arr, i = 0):
    #  dp[i] stores whether an index in the array arr can be reached via jumping
    dp = [False] * len(arr)
    dp[0] = True
    for i in range(len(arr)):
        #  If an index could not be reached from it's previous indexes then you cant go forward any more
        if not dp[i]:
            return False
        #  If you are at an index and it has a value which when added to that index exceeds the length of the arr
        #  then there's no need to search any more. You can jump to the end. So, short circuit and return true
        elif i+arr[i] >= len(arr)-1:
            return True
        #  Otherwise from an index get the max jump value and mark all the next indexes that you can jump to as true.
        else:
            for j in range(1,arr[i]+1):
                if i+j < len(arr):
                    dp[i+j] = True
    return dp[len(arr)-1]


#  This solution is more elegant. It explores the greedy method which we discarded above but with a twist.
#  For each index in the array calculate the max index that is reachable.
#  If at any point this value is less than the already existing value of max reachable index then dont do anything,
#  otherwise update the value of max reachable index.
#  If at any point, the value of max reachable index is lower than the value of the current index,
#  then return false because that indicates that this index cant be reached.
#  Time complexity of this solution is O(n)
def canJump03(arr, i = 0):
    maxReachableIndex = 0
    for i in range(len(arr)):
        if i > maxReachableIndex:
            return False
        reachableIndex = i + arr[i]
        if reachableIndex > maxReachableIndex:
            maxReachableIndex = reachableIndex
    return maxReachableIndex >= len(arr)-1



#  One more way to explore this is to think of each index of the given array as nodes in a graph
#  and the value of the element at that index as the edges from that index to the next indexes.
#  So, for example the array : [3,0,2,0,2,1,4,3]
#  can be represented as the graph :
#  {"0": [1,2,3], "1": [], "2": [3,4], "3": [], "4": [5,6], "5": [6], "6": [7], "7": []}
#  And then find a path from node 0 to node len(arr)-1 via DFS or BFS
#  Time complexity is O(n^2)
def canJump04(arr):
    if len(arr) <= 1:
        return True

    # First build the graph as a adjacency list
    graph = {}
    for i in range(len(arr)):
        adjList = []
        for j in range(1,arr[i]+1):
            if (i+j) < len(arr):
                adjList.append(i+j)
        graph[str(i)] = adjList

    # Perform BFS on the graph
    queue = [0]
    visited = {0}
    while len(queue) > 0:
        n = queue.pop(0)
        for v in graph[str(n)]:
            #  If you found the last element of the array as neighbour of a node that was reachable then there
            #  exists a path to the last index of the array. So, return True
            if v == len(arr)-1:
                return True
            else:
                if v not in visited:
                    queue.append(v)
                    visited.add(v)

    return False
