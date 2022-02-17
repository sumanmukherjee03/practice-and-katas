def describe():
    desc = """
Problem : Given a set of positive numbers, determine if a subset exists whose sum is equal to a given number S.
Example :
    Input: {1, 2, 3, 7}, S=6
    Output: True
    The given set has a subset whose sum is 6: {1, 2, 3}

    Input: {1, 2, 7, 1, 5}, S=10
    Output: True
    The given set has a subset whose sum is 10: {1, 2, 7}

    Input: {1, 3, 4, 8}, S=6
    Output: False
    The given set does not have any subset whose sum is equal to 6.

-------------
    """
    print(desc)

#  Time complexity of this solution is O(n * s)
def can_partition(nums, sum):
    n = len(nums)
    if n == 0:
        return False

    #  Initially default the dp array to -1 for each cell
    #  In dp, each cell dp[i][ts] stores the information if target sum ts can be reached by considering elements upto index i.
    #  Whether that is including element at index i or not is not that important.
    #  As long as the subproblem has considered the elements from 0 to index i it is fine.
    dp = [[-1 for x in range(0, sum+1)] for y in range(0,n)]

    #  For all elements if the target sum is 0, then the cell value will be 1
    #  because you can always have that sum with an empty subset
    for i in range(0, n):
        dp[i][0] = 1

    #  Fill in the cell value if only the first element was considered for various values of the target sum
    for s in range(1, sum+1):
        if nums[0] == s:
            dp[0][s] = 1
        else:
            dp[0][s] = 0

    for i in range(1, n):
        for s in range(1, sum+1):
            #  If upto the previous element we could achieve the target sum s, then that means we can
            #  still achieve the same with current element by not including it.
            #  If upto the previous element we could achieve the target sum (s - value of current element), then that means we can
            #  still achieve the same with current element by including it.
            if dp[i-1][s] == 1:
                dp[i][s] = 1
            elif nums[i] <= s and dp[i-1][s-nums[i]] == 1:
                dp[i][s] = 1
            else:
                dp[i][s] = 0

    return True if dp[n-1][sum] == 1 else False


#  This optimized version has the same time complexity but the dp array has a slightly lower space complexity of O(s).
#  The dp array is 1d instead of 2d and it's cells values can mutate
def can_partition_space_optimized(nums, sum):
    n = len(nums)
    if n == 0:
        return False

    dp = [False for x in range(0, sum+1)]

    #  You can always have an empty subset to get to the target sum of 0
    dp[0] = True

    #  If you only consider the first element of the nums array in the subset, and if you can reach
    #  a target sum of s within the range of 1 .. sum+1 then that implies that you can have a subset with just that 1 element
    for s in range(1, sum+1):
        if s == nums[0]:
            dp[s] = True

    #  Now process the subsets of all the numbers starting with considering numbers from index 1 onwards
    #  And consider the subset sum from sum backwards to 0 for each of those subsets of numbers
    for i in range(1, n):
        for s in range(sum, -1, -1):
            #  When considering numbers upto index i, we can reach a sum of s if we either consider the number at index i or not
            #  If we arent gonna include the number at index i in our subset, that means dp[s]'s value will remain unchanged from whatever it was if the numbers upto i-1 were considered.
            #  If we are gonna include the number at index i, then first check that target sum s is not already reached and that the required target sum is higher than the number at index i
            #  Only in that case if we include the number at i, then the value of dp[s] will be whatever the current value of dp[s-that numbers value is] because
            #  that indicates if the target sum 's-that numbers value is' could be reached by considering the elements upto i-1.
            if not dp[s] and s >= nums[i]:
                dp[s] = dp[s-nums[i]]

    return dp[sum]


def main():
    describe()

    input = [1, 2, 3, 7]
    target = 6
    print("Input : " + str(input) + ", " + str(target))
    print("Can partition: " + str(can_partition_space_optimized(input, target)))

    input = [1, 2, 7, 1, 5]
    target = 10
    print("Input : " + str(input) + ", " + str(target))
    print("Can partition: " + str(can_partition_space_optimized(input, target)))

    input = [1, 3, 4, 8]
    target = 6
    print("Input : " + str(input) + ", " + str(target))
    print("Can partition: " + str(can_partition_space_optimized(input, target)))

main()

