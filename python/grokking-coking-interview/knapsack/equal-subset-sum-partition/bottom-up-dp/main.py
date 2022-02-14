def describe():
    desc = """
Problem : Given a set of positive numbers, find if we can partition it into 2 subsets such that the sum of elements in both subsets is equal.
Example :
    Input: {1, 2, 3, 4}
    Output: True
    Explanation: The given set can be partitioned into two subsets with equal sum: {1, 4} & {2, 3}

    Input: {1, 1, 3, 4, 7}
    Output: True
    Explanation: The given set can be partitioned into two subsets with equal sum: {1, 3, 4} & {1, 7}


--------------
    """
    print(desc)

#  Time complexity of the recursive algo is O(n * s) where n is the total number of numbers and s is the sum we want to find
def can_partition(nums):
    #  If we partition the given set of numbers into 2 sets such that their sums are equal
    #  then the sum of those sets must be half of the sum of the all the numbers in nums, ie s/2.
    #  If the total sum is an odd number then the numbers cant be partitioned.
    #  From that point onwards, this problem becomes similar to the 0/1 knapsack problem,
    #  because it is a matter of finding a subset of numbers such that their total sum is s/2.
    s = sum(nums)
    if s % 2 != 0 or len(nums) == 0:
        return False

    targetSum = s // 2

    # initialize the 'dp' array, -1 for default, 1 for true and 0 for false
    #  Each cell dp[i][ts] in the dp array represents if the target sum ts can be reached by considering elements upto index i
    dp = [[-1 for x in range(0, targetSum + 1)] for y in range(0, len(nums))]

    # For target sum of 0, dont include any of the numbers
    for i in range(0, len(nums)):
        dp[i][0] = 0

    # For any target sum, if you are only considering the first element, then find the dp entries
    for ts in range(0, targetSum+1):
        if nums[0] == ts:
            dp[0][ts] = 1
        else:
            dp[0][ts] = 0

    #  Now you are filling up the dp array by considering elements upto index i and target sum of ts
    #  ie, can we reach target sum ts by considering elements upto i
    for i in range(1, len(nums)):
        for ts in range(1, targetSum+1):
            #  Now, we decide if the target sum ts can be reached by considering elements upto index i
            #  by either including element at index i or not including it.
            #  Either ways, can we reach the target sum or not.
            if dp[i-1][ts] == 1 or dp[i-1][ts-nums[i]] == 1:
                dp[i][ts] = 1
            else:
                dp[i][ts] = 0

    #  The result in a bottom up approach is the last element because that means you have considered all the elements and the final target sum you want to reach
    return True if dp[len(nums)-1][targetSum] == 1 else False


def main():
    describe()

    input = [1, 2, 3, 4]
    print("Input : " + str(input))
    print("Can partition: " + str(can_partition(input)))
    print("------------------")

    input = [1, 1, 3, 4, 7]
    print("Input : " + str(input))
    print("Can partition: " + str(can_partition(input)))
    print("------------------")

    input = [2, 3, 4, 6]
    print("Input : " + str(input))
    print("Can partition: " + str(can_partition(input)))
    print("------------------")

main()
