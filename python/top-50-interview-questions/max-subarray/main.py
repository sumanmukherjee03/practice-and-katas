#  Problem : Given a non-empty array of integers, create a function that
#  returns the sum of the subarray that has the greatest sum
#  Ex : 2,3,-6,4,2,-8,3

# This first solution is brute force and involves counting the sum of every possible subarray.
#  We cant sort because we loose the order.
#  By brute force, for every set of indexes (i,j), we can count the sum from 0...j, 1...j, 2...j and so on
#  where i keeps increasing from 0 to j for a particular j.
#  And find the max. This means 3 loops.
#  The time complexity here is O(n^3)
def maximumSubarray01(arr):
    if len(arr) == 0:
        return 0
    max = float('-inf')
    for j in range(len(arr)):
        #  Iterate till j+1 because we want to include the element at index j
        for i in range(j+1):
            sum = 0
            #  Iterate till j+1 because we want to include the element at index j
            for k in range(i, j+1):
                sum += arr[k]
            if sum > max:
                max = sum
    return max


#  An improved solution is of time complexity O(n^2)
#  We maintain 2 pointers - one moves the index i
#  And the other moves from index i to the end
#  So, calculating i...3, i...4, i...5 --- i...end
#  And at every step comparing the sum with the max
def maximumSubarray02(arr):
    if len(arr) == 0:
        return 0
    max = float('-inf')
    for i in range(len(arr)):
        sum = 0
        #  Iterate till j+1 because we want to include the element at index j
        for j in range(i,len(arr)):
            sum += arr[j]
            if sum > max:
                max = sum
    return max


#  Another slightly different approach is via dynamic programming.
#  A relationship like this somewhat simplifies the problem
#  MaxSum(4,4) = arr[4]
#  MaxSum(3,4) = max(arr[3], arr[3] + MaxSum(4,4))
#  MaxSum(1,4) = max(arr[2], arr[2] + MaxSum(3,4))
#  MaxSum(1,4) = max(arr[1], arr[1] + MaxSum(2,4))
def maximumSubarray03(arr):
    if len(arr) == 0:
        return 0
    def init2DArray(size):
        twod = []
        new = []
        for i in range (0, size):
            for j in range (0, size):
                new.append(0)
            twod.append(new)
            new = []
        return twod

    maxBetweenIndices = init2DArray(len(arr)) # This stores the max between (i, j) or the local maximum between ranges
    globalMax = arr[0]

    for j in range(1, len(arr)):
        for i in range(j,-1,-1):
            if i == j:
                maxBetweenIndices[j][j] = arr[j]
            else:
                maxBetweenIndices[i][j] = arr[i] + maxBetweenIndices[i+1][j], arr[i]
            globalMax = max(maxBetweenIndices[i][j], globalMax)
    return globalMax



#  To reduce the time complexity further we can make use of Kadane's algo
#  You can read about the algo here : https://medium.com/@rsinghal757/kadanes-algorithm-dynamic-programming-how-and-why-does-it-work-3fd8849ed73d
#  A relationship like this simplifies the problem significantly.
#  MaxSumAtIndex(4) = max(arr[4], arr[4] + MaxSumAtIndex(3))
#  MaxSumAtIndex(3) = max(arr[3], arr[3] + MaxSumAtIndex(2))
#    Formula : local_max[i] = max(arr[i], arr[i] + local_max[i-1])
#  In layman's terms max subarray is either
#    - sum of the previous subarray with highest sum + the current number
#    - OR the current number
#        The 2 conditions combined above gurantees the contiguosness of the subarray
def maximumSubarray04(arr):
    globalMax = float('-inf')
    localMax = 0
    for elm in arr:
        localMax = max(elm, localMax + elm)
        globalMax = max(localMax, globalMax)
    return globalMax


#  To represent the same problem above but in a dynamic programming construct
def maximumSubarray05(arr):
    globalMax = float('-inf')
    dp = []
    dp[0] = arr[0]
    for i in range(1, len(arr)):
        dp[i] = max(arr[i], arr[i] + dp[i-1])
        globalMax = max(dp[i], globalMax)
    return globalMax
