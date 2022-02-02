def describe():
    desc = """
Problem : Given an array of numbers sorted in an ascending order, find the ceiling of a given number key.
The ceiling of the key will be the smallest element in the given array greater than or equal to the key.

Example :
    Input: [4, 6, 10], key = 6
    Output: 1
    Explanation: The smallest number greater than or equal to 6 is 6 having index 1.

    Input: [1, 3, 8, 10, 15], key = 12
    Output: 4
    Explanation: The smallest number greater than or equal to 12 is 15 having index 4.

-----------
    """
    print(desc)

#  This is the same as order agnostic binary search. Only difference being that if we cant find the element
#  then the next big number will be pointed out by the element at index start.

def search_ceiling_of_a_number(arr, key):
    if len(arr) == 0:
        return -1

    start, end = 0, len(arr) - 1
    while start <= end:
        mid = start + (end-start)//2
        if arr[mid] == key:
            return mid
        if key < arr[mid]:
            end = mid - 1
        else:
            start = mid + 1
    return start


def main():
    describe()

    input = [4, 6, 10]
    key = 6
    print("Input : " + str(input) + " , " + str(key))
    print(search_ceiling_of_a_number(input, key))

    input = [1, 3, 8, 10, 15]
    key = 12
    print("Input : " + str(input) + " , " + str(key))
    print(search_ceiling_of_a_number(input, key))

    input = [4, 6, 10]
    key = 17
    print("Input : " + str(input) + " , " + str(key))
    print(search_ceiling_of_a_number(input, key))

    input = [4, 6, 10]
    key = -1
    print("Input : " + str(input) + " , " + str(key))
    print(search_ceiling_of_a_number(input, key))

main()
