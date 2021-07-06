#  Problem : Given a string of alphabetical chars only create a function that returns the length
#  of the longest possible palindrome from the characters of the string
#  Ex: Input - "abbaba"
#  Output : 5 - Longest palindrome possible is "baaab"

#  Building a palindrome involves
#       - if there is already a palindrome then adding the same char at the front and back will also produce a palindrome.
#  Remember we are just asked to find the length of the longest palindrome, not the longest palindrome itself.
#  As long as some char occurs in even number in the string we can make a palindrome out of it.
#  If a char occurs in odd number we can still use (num of occurances - 1) of that char in building the longest palindrome.
#  Also for chars with odd number of occurances, if they exist in the string we can only one make use of 1 of those such char in the palindrome.
#       Only in the middle of the palindrome.
# The time complexity is O(n)
def longestPalindrome(str):
    visited = {}

    for c in str:
        if visited.get(c):
            visited[c] += 1
        else:
            visited[c] = 1

    maxLen = 0
    oddOccurancesCharCount = 0
    for k in visited:
        if visited[k] % 2 == 0:
            maxLen += visited[k]
        else:
            maxLen += visited[k] - 1
            oddOccurancesCharCount += 1

    if oddOccurancesCharCount > 0:
        maxLen += 1

    return maxLen

