def describe():
    desc = """
Problem : Given a string, find the length of the longest substring in it with no more than K distinct characters.
For example :
    Input: String="araaci", K=2
    Output: 4
    Explanation: The longest substring with no more than '2' distinct characters is "araa", where the distinct chars are 'a' and 'r'.

-----------
    """
    print(desc)


#  Time complexity is O(n)
def find_substr_with_distinct_chars(str, k):
    if len(str) == 0 or k <= 0:
        return ""
    maxSubstr = "" # Maintains final result
    charMap = {}
    head = 0 # Pointer of the start of the sliding window
    # tail is the end pointer of the sliding window
    for tail in range(0, len(str)):
        tailChar = str[tail]
        if tailChar not in charMap:
            charMap[tailChar] = 0
        charMap[tailChar] += 1

        while len(charMap) > k:
            headChar = str[head]
            charMap[headChar] -= 1
            if charMap[headChar] == 0:
                del charMap[headChar]
            head += 1

        substr = str[head:tail+1]
        if len(substr) > len(maxSubstr):
            maxSubstr = substr

    return maxSubstr


def main():
    describe()
    str = "araaci"
    k = 2
    res = find_substr_with_distinct_chars(str, k)
    print("Input", str, k)
    print("Longest substring with k distinct chars is : ", res)
    print("Length of longest such substring is : ", len(res))

main()
