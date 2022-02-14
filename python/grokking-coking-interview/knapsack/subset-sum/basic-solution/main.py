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
    return can_partition_recursive(nums, sum)

#  Time complexity of this solution is O(2^n) because there are 2 choices, either to include a current number or not
def can_partition_recursive(nums, currentSum, currentIndex = 0):
    if currentSum == 0 and currentIndex <= len(nums):
        return True
    if len(nums) == 0 or currentIndex >= len(nums):
        return False
    #  If the current number is lower than the sum then try and include that number and check if the target sum can be reached with the next numbers
    #  But if the number is larger than the current target sum or the target sum could not be reached by including this numbers
    #  then dont include this number but carry on with the next iteration
    if nums[currentIndex] <= currentSum:
        if can_partition_recursive(nums, currentSum-nums[currentIndex], currentIndex+1):
            return True
    return can_partition_recursive(nums, currentSum, currentIndex+1)

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
