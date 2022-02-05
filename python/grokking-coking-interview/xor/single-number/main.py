def describe():
    desc = """
Problem : In a non-empty array of integers, every number appears twice except for one, find that single number.
Example:
    Input: 1, 4, 2, 1, 3, 2, 3
    Output: 4


--------------
    """
    print(desc)

#  One easy solution is that we maintain a hashmap with keys as the numbers and value as a boolean
#  If you encounter the number second time, remove it from the hash map. That way at the end you will be left with only 1 number.
#  A better solution is using xor.
#  Since each number is repeated twice except one, most of them will cancel out when xor is applied.
#  Time complexity is O(n)
def find_single_number(arr):
    x = arr[0]
    for i in range(1, len(arr)):
        x = x ^ arr[i]
    return x

def main():
    describe()
    arr = [1, 4, 2, 1, 3, 2, 3]
    print("Input : " + str(arr))
    print("Output : " + str(find_single_number(arr)))

main()
