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
    return recursive_partition(nums)

def recursive_partition(nums, set1 = [], set2 = [], index = 0):
    s1 = set1[:]
    s2 = set2[:]
    s1.append(nums[index])
    s2.append(nums[index])
    if index == len(nums)-1:
        if sum(s1) == sum(set2):
            return True
        elif sum(set1) == sum(s2):
            return True
        else:
            return False
    else:
        return recursive_partition(nums, s1, set2, index+1) or recursive_partition(nums, set1, s2, index+1)


def main():
    describe()

    input = [1, 2, 3, 4]
    print("Input : " + str(input))
    print("Can partition: " + str(can_partition(input)))

    input = [1, 1, 3, 4, 7]
    print("Input : " + str(input))
    print("Can partition: " + str(can_partition(input)))

    input = [2, 3, 4, 6]
    print("Input : " + str(input))
    print("Can partition: " + str(can_partition(input)))

main()
