def describe():
    desc = """
Problem : We are given an array containing n objects.
Each object, when created, was assigned a unique number from the range 1 to n based on their creation sequence.
This means that the object with sequence number 3 was created just before the object with sequence number 4.
Write a function to sort the objects in-place on their creation sequence number in O(n)O(n) and without using any extra space.

Example :
    Input: [3, 1, 5, 4, 2]
    Output: [1, 2, 3, 4, 5]

    Input: [1, 5, 6, 4, 3, 2]
    Output: [1, 2, 3, 4, 5, 6]
-----------------

    """
    print(desc)

#  We iterate the array one number at a time, and if the current number we are iterating is not at the correct index, we swap it with the number at its correct index.
#  Time complexity of this sort is O(n)
def cyclic_sort(nums):
    i = 0
    while i < len(nums):
        j = nums[i]
        if i+1 == j:
            i += 1
        else:
            nums[j-1], nums[i] = j, nums[j-1]
    return nums


def main():
    describe()
    arr = [1, 5, 6, 4, 3, 2]
    print("Input : ", arr)
    cyclic_sort(arr)
    print("Output : ", arr)

main()
