#  Problem : Given an array of integers create a function that returns the index of a peak element.
#  A peak element is one which is greater than or equal to it's neighbours, ie the next and the previous element.
#  Assume that arr[-1] and arr[n] are -infinity.
#  If there are multiple peaks return the index of one of them.

# This is the simplest solution and has O(n) time complexity
def findPeak01(arr):
    if len(arr) == 0:
        return -1
    if len(arr) == 1:
        return 0
    for i in range(len(arr)):
        if i+1 < len(arr) and i-1 >= 0:
            if arr[i] >= arr[i+1] and arr[i] >= arr[i-1]:
                return i
        elif i+1 < len(arr) and i-1 < 0:
            if arr[i] >= arr[i+1]:
                return i
        elif i+1 >= len(arr) and i-1 >= 0:
            if arr[i] >= arr[i-1]:
                return i
    return -1



# This solution is based on divide and conquer.
# Find the mid element.
# If the mid element is lower than mid+1 it means that there is a upward trend on the right of mid
#     in which case a peak will exist on the right
#         regardless of whether there is a flatline on the right after an upward climb or not
#          ____
#         /           /\           /
#       /           /   \        /
#     /           /            /
# If the mid element is higher than mid+1 it means there is an upward trend on the left of mid+1
#     in which case a peak will exist on the left
#         regardless of whether there is a flatline on the left after an upward climb
# If the mid element is at the same level as the mid+1 then 3 scenarios are possible on the left
#     - all elements on the left are of the same height in which case there is a peak
#     - mid is the same as mid+1 and then there is an upward trend on left of mid in which case there is a peak
#     - mid is the same as mid+1 and then there is a downward trend on the left in which case mid is the peak
#      ___        /\                        /            /
#          \     /   \                  ___/            /
#            \         \               /         _____ /
#             \          \ ___        /

# This is similar to a binary search and the time complexity is O(logn)
def findPeak02(arr):
    left = 0
    right = len(arr)-1
    while left < right:
        mid = (left+right)/2
        if arr[mid] < arr[mid+1]:
            left = mid+1
        else:
            right = mid
    return left
