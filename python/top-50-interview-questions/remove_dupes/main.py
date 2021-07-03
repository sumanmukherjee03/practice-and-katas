# Problem : Given an array of integers create a fucntion that returns an array
#  containing the values of the array without duplicates - order doesnt matter

#  One solution is to traverse the array and push the element into an output array if it already doesnt exist in the output array.
#  The time complexity of this is O(n^2)
def removeDuplicates01(arr):
    res = []
    for elm in arr:
        if elm not in res:
            res.append(elm)
    return res


# Sort the array and then traverse it. If element does not exist in the resulting array's last element
#  then insert it into the resulting element.
#  Time complexity is O(nlogn) for sort and O(n) for traversal -> O(nlogn)
def removeDuplicates02(arr):
    arr.sort()
    res = []
    for i in range(len(arr)):
        if len(res) == 0 or res[len(res)-1] != arr[i]:
            res.append(arr[i])
    return res


# Same solution as above, but with a slightly different implementation
def removeDuplicates03(arr):
    if len(arr) == 0:
        return []
    arr.sort()
    res = [arr[0]]
    for i in range(1, len(arr)):
        if arr[i] != arr[i-1]:
            res.append(arr[i])
    return res


# This solution leverages the fact that if you insert the same key in a hash table with some value it will simply get overriden
#  The time complexity of this is O(n)
def removeDuplicates04(arr):
    visited = {}
    for elm in arr:
        visited[elm] = True
    return list(visited.keys())
