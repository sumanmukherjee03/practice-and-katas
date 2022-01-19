def describe():
    desc = """
Problem: Given an array of sorted numbers, remove all duplicates from it. You should not use any extra space.
After removing the duplicates in-place return the length of the subarray that has no duplicate in it.
Example:
    Input: [2, 3, 3, 3, 6, 9, 9]
    Output: 4
    Explanation: The first four elements after removing the duplicates will be [2, 3, 6, 9].

--------------
    """
    print(desc)

#  As the input array is sorted, therefore, one way to do this is to shift the elements left
#  whenever we encounter duplicates. In other words, we will keep one pointer for iterating the array
#  and one pointer for placing the next non-duplicate number.
#  So our algorithm will be to iterate the array and whenever we see a non-duplicate number we move it next to the last non-duplicate number we’ve seen.
#  Time complexity O(n) and space complexity is O(1)
def remove_dupes(arr):
    head = 0 # Represents the next non dupe element
    tail = 0 # Pointer for iteration
    for tail in range(1, len(arr)):
        if arr[head] != arr[tail]:
            arr[head+1], arr[tail] = arr[tail], arr[head+1]
            head += 1
    return head


def main():
    describe()
    arr = [2, 3, 3, 3, 6, 9, 9]
    print("Input : ", arr)
    index = remove_dupes(arr)
    print("Rearranged array : ", arr)
    print("Resulting sub array : ", arr[0:index+1])

main()
