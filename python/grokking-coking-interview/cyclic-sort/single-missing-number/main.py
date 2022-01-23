def describe():
    desc = """
Problem : We are given an array containing n distinct numbers taken from the range 0 to n.
Since the array has only n numbers out of the total n+1 numbers, find the missing number.

Example:
    Input: [4, 0, 3, 1]
    Output: 2

    Input: [8, 3, 5, 2, 4, 6, 0, 1]
    Output: 7
------------------
    """
    print(desc)


def find_missing_number(nums):
    i = 0

    #  If the number is larger than the index, ie "n", then ignore it (ie dont swap) and move to the next one.
    #  At the end of the iteration that number n, will be the one sitting in the wrong place.
    #  Rest of the algo is the same as cyclic sort.
    while i < len(nums):
        j = nums[i]
        if j < len(nums) and i != j:
            nums[j], nums[i] = j, nums[j]
        else:
            i += 1

    for k in range (0, len(nums)):
        if k != nums[k]:
            break

    return k



def main():
    describe()
    nums = [8, 3, 5, 2, 4, 6, 0, 1]
    print("Input :", nums)
    res = find_missing_number(nums)
    print("Result :", res)

main()
