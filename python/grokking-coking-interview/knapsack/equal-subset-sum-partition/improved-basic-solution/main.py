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

#  Time complexity of the recursive algo is O(2^n)
def can_partition(nums):
    #  If we partition the given set of numbers into 2 sets such that their sums are equal
    #  then the sum of those sets must be half of the sum of the all the numbers in nums, ie s/2.
    #  If the total sum is an odd number then the numbers cant be partitioned.
    #  From that point onwards, this problem becomes similar to the 0/1 knapsack problem,
    #  because it is a matter of finding a subset of numbers such that their total sum is s/2.
    s = sum(nums)
    if s % 2 != 0:
        return False
    return recursive_partition(nums, s/2, 0)

#  Time complexity is O(2^n) because at each iteration there are 2 choices, either include the current element in set 1 or set 2 and there are only 2 sets.
def recursive_partition(nums, sum, currentIndex):
    #  If the remaining sum you are trying to find is down to 0, then you have already found all the numbers for 1 subset
    if sum == 0:
        return True
    if len(nums) == 0 or currentIndex >= len(nums):
        return False
    #  If the current number is less than the remaining sum, then there are only 2 choices, either include it or not include it
    if nums[currentIndex] <= sum:
        if recursive_partition(nums, sum - nums[currentIndex], currentIndex+1):
            return True
    return recursive_partition(nums, sum, currentIndex+1)


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
