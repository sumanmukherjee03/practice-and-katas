def describe():
    desc = """
Problem: Given an array of sorted numbers and a target sum, find a pair in the array whose sum is equal to the given target.
Example:
    Input: [1, 2, 3, 4, 6], target=6
    Output: [1, 3]
    Explanation: The numbers at index 1 and 3 add up to 6: 2+4=6

------------
    """
    print(desc)

#  Time complexity is O(n)
def two_sum(arr, k):
    head = 0
    tail = len(arr) - 1
    while head < tail:
        s = arr[head] + arr[tail]
        if s == k:
            return [head, tail]
        elif s < k:
            head += 1
        else:
            tail -= 1
    return [-1, -1]

def main():
    describe()
    arr = [1,2,3,4,6]
    k = 6
    res = two_sum(arr, k)
    print("Input : ", arr, k)
    print("Resulting indices are : ", res)

main()
