# Problem : Given an array of integers of size at least n+1 and n elements to fill it,
#  create a function to return the duplicate element. Based on pigeonhole principle there will be at least 1 duplicate element.
#  Assume that there is only 1 duplicate and that the duplicate element can be repeated many times

#  For each element in the array traverse the remaining elements and see if you find a duplicate element.
#  The time complexity of this brute force solution is O(n^2)
def findDuplicate01(arr):
    for i in range(len(arr)):
        for j in range(i+1,len(arr)):
            if arr[i] == arr[j]:
                return arr[i]
    return None

#  Another solution is to sort the array first and then return the repeated element.
#  That algo will have complexity of O(nlogn)

# This one has time complexity of O(n)
def findDuplicate02(arr):
    visited = {}
    for elm in arr:
        if visited.get(elm):
            return elm
        else:
            visited[elm] = True
    return None


# TODO : We can come back to this solution afterwards.
#
#  An even better solution is to solve it with Floyd Cycle Detection algorithm
#  An important information to keep in mind that in a single link linked list
#  if a node has 2 next pointers pointing to it, then that's the one which is causing the cycle in the linked list.
#  This is also known as the tortoise and hare technique.
#  It involves 2 pointers - a slow pointer and a fast pointer.
