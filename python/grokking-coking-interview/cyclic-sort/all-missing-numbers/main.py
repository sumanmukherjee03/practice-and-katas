def describe():
    desc = """
Problem : We are given an unsorted array containing numbers taken from the range 1 to n.
The array can have duplicates, which means some numbers will be missing. Find all those missing numbers.

Example :
    Input: [2, 3, 1, 8, 2, 3, 5, 1]
    Output: 4, 6, 7
    Explanation: The array should have all numbers from 1 to 8, due to duplicates 4, 6, and 7 are missing.
------------------


    """
    print(desc)


#  Time complexity is O(n)
def find_missing_numbers(nums):
    i = 0

    #  This problem is the same as cyclic sort except for the fact that some of the numbers are repeated
    #  and while swapping you will notice that there are numbers in the correct position already.
    while i < len(nums):
        j = nums[i]
        if i != nums[i-1]-1 and j-1 != nums[j-1]-1:
            nums[j-1], nums[i] = nums[i], nums[j-1]
        else:
            i += 1

    out = []
    for k in range (0, len(nums)):
        if k != nums[k]-1:
            out.append(k+1)

    return out



def main():
    describe()
    nums = [2, 3, 1, 8, 2, 3, 5, 1]
    print("Input :", nums)
    res = find_missing_numbers(nums)
    print("Result :", res)

main()

