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
    val, set1, set2 = recursive_partition(nums)
    print("solution : " + str(set1) + ", " + str(set2))
    return val

#  Time complexity is O(2^n) because at each iteration there are 2 choices, either include the current element in set 1 or set 2 and there are only 2 sets.
def recursive_partition(nums, set1 = [], set2 = [], index = 0):
    s1 = set1[:] # copy of set 1, because we dont want to modify the original set
    s2 = set2[:] # copy of set 2, because we dont want to modify the original set
    s1.append(nums[index]) # include current number in copy of set 1 so that we can check against unaltered set 2
    s2.append(nums[index]) # include current number in copy of set 2 so that we can check against unaltered set 1
    #  We are considering 1 index at a time and if this is the last index we are considering
    #  then check if the sum of the 2 subsets are equal.
    #  When doing so, first include the current element in set 1 and then check and afterwards if that doesnt produce the desired result
    #  then include the current element in the set 2 and then check
    if index == len(nums)-1:
        if sum(s1) == sum(set2):
            return True, s1, set2
        elif sum(set1) == sum(s2):
            return True, set1, s2
        else:
            return False, [], []
    else:
        val1, set11, set12 = recursive_partition(nums, s1, set2, index+1)
        val2, set21, set22 = recursive_partition(nums, set1, s2, index+1)
        if val1:
            return val1, set11, set12
        elif val2:
            return val2, set21, set22
        else:
            return False, [], []


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
