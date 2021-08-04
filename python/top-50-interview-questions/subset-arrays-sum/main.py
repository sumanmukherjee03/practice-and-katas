#  Problem : Given an array of strictly positive integers and an integer k, create a function
#  that returns the number of subsets of arr that sum upto k.
#  Ex: Input - [1,2,3,1], k = 4
#  Output: 3
#  Explanation : [1,2,1] ; [1,3] ; [3,1]

#  In an array of n elements we can have (2^n) subsets - why? because there are 2 possibilities for each element - to take it or leave it.
#  One details to notice here is that the elements are strictly positive, so local sum can only increase.
#  Also, we dont need the actual subsets, we just need the count of the number of subsets.
#  At each index there can be 2 choices, either take an element of the array or leave it.
#  i=0                             []
#                                  /\
#                                /    \
#  i=1                         [1]    []
#  Using this concept recursively, if going down a path you encounter a sum of k, then return back 1 because you found 1 path to get there and backtrack.
#  If sum found is greater than k then return 0 and backtrack because everything down from that path will be more than k.

#  Here i is the index in the array you are considering
#  And s is the local sum you are considering to check if it has reached k or not
#  Time complexity is O(2^n)
def subsetsThatSumUpToK(arr, k, i = 0, s = 0):
    if s == k:
        return 1
    if s > k or i >= len(arr):
        return 0
    return subsetsThatSumUpToK(arr, k, i+1, s + arr[i]) + subsetsThatSumUpToK(arr, k, i+1, s)
