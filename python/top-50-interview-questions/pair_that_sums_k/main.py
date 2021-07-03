#  Problem : Given an array and a number k, find out of there are 2 numbers that add upto k

#  Solution 1 : This is the brute force solution
# Time complexity of this solution is O(n^2) because of the nested loop
def findPair01(arr, k):
    for i in range(len(arr)):
        for j in range(i+1, len(arr)):
            if arr[i]+arr[j] == k:
                return True
    return False


# Solution 2 : We first sort the array because it is easier to operate on a sorted array.
#  This uses 2 pointers - one left and one right
#  Since the array is sorted, moving the left forward pointer will increase the sum of left+right
#  And moving the right backward will reduce the sum of left+right
# Time complexity is O(nlogn) because sorting is nlogn and time for traversing is n.
def findPair02(arr, k):
    arr = arr.sort()
    left = 0
    right = len(arr)-1
    while left > right:
        if arr[left]+arr[right] == k:
            return True
        if arr[left]+arr[right] < k:
            left += 1
        if arr[left]+arr[right] > k:
            right -= 1
    return False


# Solution 3 : We use a hash table to store the remainder of the sum for each element that we visit
#  If a number exists with that value we found what we needed and return true
#  This has a time complexity of O(n) because hash table lookup and insertion is O(1)
def findPair03(arr, k):
    visited = {}
    for elm in arr:
        if visited.get(k-elm):
            return True
        else:
            visited[k-elm] = True
    return False
