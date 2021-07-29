#  Problem : Given a non-empty rotated sorted array of integers that has no duplicates,
#  create a function that returns the minimum.
#  The array is sorted in ascending order and that it's rotated by an unknown number of positions to the right.
#  Ex : Input - [4,5,1,2,3]
#  Output - 1

def minimum(arr):
    if len(arr) == 1:
        return arr[0]

    #  Starting from index 1 find the point in the array where the previous element is greater than the current element
    i = 1
    prev = arr[0]
    while i < len(arr):
        current = arr[i]
        if current < prev:
            break
        prev = current
        i += 1
    # If you travered to the end of the array but did not encounter a break, it means the array was sorted properly
    # in which case return the first element, otherwise return the last current found
    if i < len(arr):
        return current
    else:
        return arr[0]


#  This second approach is based on divide and conquer.
#  Imagine the situation like a graph :
#          /\    /
#        /    \/
#      /
#  Consider that you calculate the mid pointer always.
#  We want to find the mid pointer such that we finally find the breaking element, so that we can return the next one as the minimum
#  4 conditions :
#       - if the middle element is the element at the peak in the middle, then the minimum is the element at mid+1
#       - if the middle element is the element at the dip in the middle, then the minimum is the element at mid
#       - if the middle element is somewhere where the peak and dip are not there
#           - then is the element at the end smaller than the mid?
#               - if yes, then search in the right half, ie mid+1 to end
#               - otherwise, search in the left half, ie start to mid-1
#  Time complexity is O(nlogn)
def minimum(arr):
    left = 0
    right = len(arr)-1
    #  If the first element is smaller than the last element,
    #  it indicates that the array has been rotated few times
    #  but it is back to it's original initial sorted state
    if arr[left] < arr[right]:
        return arr[left]
    while left < right:
        mid = (left+right)/2
        # This condition indicates that we found the break, so return the element right after it
        if arr[mid+1] < arr[mid]:
            return arr[mid+1]
        elif arr[mid] < arr[mid-1]:
            return arr[mid]
        elif arr[mid] > arr[right]:
            left = mid+1
        else:
            right = mid-1
    return arr[left]


# Same solution as above but recursively
def minimum(arr):
    def minRec(arr, leftPtr, rightPtr):
        if leftPtr >= rightPtr:
            return arr[leftPtr]
        if arr[leftPtr] < arr[rightPtr]:
            return arr[leftPtr]
        mid = (left+right)/2
        if arr[mid+1] < arr[mid]:
            return arr[mid+1]
        elif arr[mid] < arr[mid-1]:
            return arr[mid]
        elif arr[mid] > arr[rightPtr]:
            return minRec(arr, mid+1, rightPtr)
        else:
            return minRec(arr, leftPtr, mid-1)
    return minRec(arr, 0, len(arr)-1)


