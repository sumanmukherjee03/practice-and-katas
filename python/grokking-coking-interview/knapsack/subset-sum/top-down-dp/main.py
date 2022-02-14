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

def can_partition(nums, sum):
    #  Initially default the dp array to -1 for each cell
    #  In dp, each cell dp[i][ts] stores the information if target sum can be reached by considering elements at index i
    dp = [[-1 for x in range(0, sum+1)] for y in range(0,len(nums))]
    val = can_partition_recursive(dp, nums, sum)
    return True if val == 1 else False

#  Time complexity of this solution is O(n * s) where n is the number of numbers and s is the target sum
def can_partition_recursive(dp, nums, currentSum, currentIndex = 0):
    if currentSum == 0:
        return 1
    if len(nums) == 0 or currentIndex >= len(nums):
        return 0
    if dp[currentIndex][currentSum] < 0:
        #  If the current number is lower than the sum then try and include that number and check if the target sum can be reached with the next numbers
        #  But if the number is larger than the current target sum or the target sum could not be reached by including this numbers
        #  then dont include this number but carry on with the next iteration
        if nums[currentIndex] <= currentSum:
            if can_partition_recursive(dp, nums, currentSum-nums[currentIndex], currentIndex+1) == 1:
                dp[currentIndex][currentSum] = 1
                return 1
        dp[currentIndex][currentSum] = can_partition_recursive(dp, nums, currentSum, currentIndex+1)
    return dp[currentIndex][currentSum]

def main():
    describe()

    input = [1, 2, 3, 7]
    target = 6
    print("Input : " + str(input) + ", " + str(target))
    print("Can partition: " + str(can_partition(input, target)))

    input = [1, 2, 7, 1, 5]
    target = 10
    print("Input : " + str(input) + ", " + str(target))
    print("Can partition: " + str(can_partition(input, target)))

    input = [1, 3, 4, 8]
    target = 6
    print("Input : " + str(input) + ", " + str(target))
    print("Can partition: " + str(can_partition(input, target)))

main()

