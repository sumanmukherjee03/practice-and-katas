def describe():
    desc = """
Problem : Given a sorted array of numbers, find if a given number 'key' is present in the array.
Though we know that the array is sorted, we don't know if it's sorted in ascending or descending order.
You should assume that the array can have duplicates.
Write a function to return the index of the 'key' if it is present in the array, otherwise return -1.

Example :
    Input: [4, 6, 10], key = 10
    Output: 2

-------------
    """
    print(desc)

def recursive_binary_search(arr, key, start = 0, stop = -1):
    #  Take care of initialization
    if stop == -1:
        stop = len(arr)

    if len(arr) == 0:
        return -1

    if stop == start:
        if arr[start] == key:
            return start
        else:
            return -1

    #  One way to calculate mid is : mid = int((start + stop)/2)
    #  However, this might result in integer overflow if both start and stop indices are very high
    # The process below is a much safer way to do the same
    mid = start + (stop-start)//2
    if key == arr[mid]:
        return mid

    idx_left, idx_right = -1, -1
    idx_left = recursive_binary_search(arr, key, start, mid)
    if mid+1 < stop:
        idx_right = recursive_binary_search(arr, key, mid+1, stop)
    if idx_left >= 0:
        return idx_left
    elif idx_right >= 0:
        return idx_right
    else:
        return -1




#  Iterative binary search
def binary_search(arr, key):
    start, end = 0, len(arr)-1
    is_ascending = arr[start] < arr[end]
    while start <= end:
        mid = start + (end - start)//2
        if key == arr[mid]:
            return mid
        if is_ascending:
            if key < arr[mid]:
                end = mid - 1
            else:
                start = mid + 1
        else:
            if key < arr[mid]:
                start = mid + 1
            else:
                end = mid - 1
    return -1



#  Since the search range is getting reduced by 1/2 on each iteration or recursive call the time complexity is O(logn)
def main():
    describe()

    input = [1, 2, 3, 4, 5, 6, 7]
    key = 5
    print("Input : " + str(input) + " , " + str(key))
    print(recursive_binary_search(input, key))

    input = [4, 6, 10]
    key = 10
    print("Input : " + str(input) + " , " + str(key))
    print(binary_search(input, key))

    input = [10, 6, 4]
    key = 10
    print("Input : " + str(input) + " , " + str(key))
    print(binary_search(input, key))

    input = [10, 6, 4]
    key = 4
    print("Input : " + str(input) + " , " + str(key))
    print(binary_search(input, key))

main()
