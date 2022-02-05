def describe():
    desc = """
Problem : Given an array of n-1n−1 integers in the range from 11 to nn, find the one number that is missing from the array.
Example :
    Input: 1, 5, 2, 6, 4
    Answer: 3

-----------
    """
    print(desc)

#  This is an easy solution. Sum all numbers from 1 to n.
#  Then keep subtracting the array elements from that sum. Whatever is left behind gives the missing number.
def crude_find_missing_number(arr):
    n = len(arr) + 1
    s1 = 0
    for num in range(1, n+1):
        s1 += num

    for i in range(0, len(arr)):
        s1 -= arr[i]

    return s1


#  This is a better solution than the previous one because it does not cause overflow in case of large numbers
#  Time complexity is O(n) and space complexity is O(1)
def find_missing_number(arr):
    n = len(arr) + 1
    x1 = 1
    for num in range(2, n+1):
        x1 = x1 ^ num

    x2 = arr[0]
    for i in range(1, len(arr)):
        x2 = x2 ^ arr[i]

    return x1 ^ x2

def main():
    describe()
    result = find_missing_number([1,5,2,6,4])
    print(str(result))

main()
