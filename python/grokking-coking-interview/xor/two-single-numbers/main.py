def describe():
    desc = """
Problem : In a non-empty array of numbers, every number appears exactly twice except two numbers that appear only once.
Find the two numbers that appear only once.
Example :
    Input: [1, 4, 2, 1, 3, 5, 6, 2, 3, 5]
    Output: [4, 6]

-----------

    """
    print(desc)

#  Time complexity is O(n)
#  Remember - a number repeated twice is a good indicator that xor is a likely solution
def find_single_numbers(arr):
    if len(arr) == 0:
        return [-1, -1]
    n1xn2 = arr[0]
    for i in range(1, len(arr)):
        n1xn2 ^= arr[i]
    #  At this point n1xn2 is the xor of the 2 single numbers n1 and n2
    #  The fact that this number is not zero means that between n1 and n2 there must be 1 bit which is 1.
    #  So, find the rightmost bit that is 1 for n1xn2
    rightmost_set_bit = 1
    while (rightmost_set_bit & n1xn2) == 0:
        #  Bit shift operator moves the set bit by 1 position to the left.
        #  So, n = 1, then n = 2, then n = 4 and so on. 1, 10, 100, 1000 in terms of binary representation
        rightmost_set_bit = rightmost_set_bit << 1

    #  At this point we have a number rightmost_set_bit such that only 1 bit is set in that.
    #  If we perform a logical AND of this number with any other number then only those numbers will be non zero
    #  that have the bit set. AND with all other numbers will be 0.
    #  This means this operation will divide the set into 2 subsets and each of those subsets will contain one of num1 and num2
    #  Finally we perform the same operation that we performed in finding the single number on each of these subsets.
    num1, num2 = 0, 0
    for i in range(0, len(arr)):
        if (arr[i] & rightmost_set_bit) != 0:
            num1 ^= arr[i]
        else:
            num2 ^= arr[i]
    return [num1, num2]

def main():
    describe()
    input = [1, 4, 2, 1, 3, 5, 6, 2, 3, 5]
    print('Input : ' + str(input))
    print('Single numbers are:' + str(find_single_numbers(input)))
    input = [2, 1, 3, 2]
    print('Input : ' + str(input))
    print('Single numbers are:' + str(find_single_numbers(input)))

main()
