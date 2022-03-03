def describe():
    desc = """
Problem : Given a string return the shortest palindrome that we can get by adding chars in front of the string.
Example :
    Input : "acbcabcbb"
    Output : "bbcbacbcabcbb"

    Input : "bcaacb"
    Output : "bcaacb"

-----------------

    """
    print(desc)

#  Take the string acbcabcbb for example. The main goal is to find the longest palindrome prefix in the string.
#  That way we can just visit the rest of chars in the string and reverse them and add it in front.
#
#  In our example string "acbcabcbb", "acbca" is the palindrome prefix.
#  So, to this string, we just need to add the rest of the chars in the reverse order.
#
#  It will always have a palindromic prefix because if there were some other chars in the front which are not palindromic, say the string was "xyzacbcabcbb"
#  ie if xyz was before the palindromic prefix then zyx needs to be inserted after the palindromic prefix.
#  ie shortest palindome from that string would have become "bbcbxyzacbcazyxbcbb". So, you see how we had to add the string zyx in the middle of the string.
#  But that is not what the problem says. The problem states we can get a palindrome by adding chars at the front of the string.
#  Hence the palindome contained in the input string would always start from the beginning, ie there will always be a palindromic prefix in this case.
#
#  The palindrome prefix can infact be treated like a mirror. And the rest of the string is reflected off the mirror.
#  This implementation is O(n^2)
def crudeShortestPalindrome(s):
    mirrorLength = 0
    #  Iterate over each char in the string and check if the substring upto that index is a palindrome.
    #  Keep incrementing the length of the mirror, ie the palindrome prefix.
    #  Trick to reverse a string in python - s[::-1]
    for i in range(len(s)+1):
        if s[0:i] == s[0:i][::-1]:
            mirrorLength = i
    return s[mirrorLength:][::-1] + s

#  A good trick to find the palindromic prefix is considering the fact that
#  if we reverse the input string the palindome would still be a substring of the reversed string.
#  ie reverse of "acbcabcbb" is "bbcbacbca" and we can see that the substring "acbca" which is the palindrome is still a substring of the reversed string.
#  The only difference being that it is present at the very end.
#  We use this knowledge about the position and join the input string with it's reversed string "acbcabcbb#bbcbacbca" separated by a separator lets say.
#  Now we can observe that the palindrome which is a prefix of the input string is also a suffix of the reverse of input string.
#  This allows us to use the longest-prefix-suffix (find the longest prefix of a string which is also a suffix) problem to find the mirror length.
#
#  Time complexity is O(n)
def shortestPalindrome(s):
    s1 = s + '#' + s[::-1]
    lpsArr = getLpsArr(s1)
    mirrorLength = lpsArr[-1]
    return s[mirrorLength:][::-1] + s

#  LPS array is longest-prefix-suffix array. This is an array where each cell represents
#  the length of the longest prefix that is also a suffix of the substring that ends at that index.
#  That way the last cell will contain the length of the longest prefix that is also suffix of the entire string.
#  Example :
#  a c b c a b c b b # b b c b a c b c a  - string
#  0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 2 3 4 5  - lpsArr
def getLpsArr(s):
    lpsArr = [0] * len(s) # Initialize with 0
    length = 0
    i = 1 # We start with index 1 because lps of 0th index is always 0
    while i < len(s):
        #  s[length] keeps track of the chars from the beginning of the string
        #  s[i] keeps track of chars in the laster parts of the string
        #  ie s[length] gives the next char from the front to match and s[i] is the next char in the iteration to match
        if s[i] == s[length]:
            length += 1
            lpsArr[i] = length
            i += 1
        #  For the next condition consider the string :
        #  a b c d a b a b c d m a
        #  0 0 0 0 1 2 1 2 3 4 0 1
        #  After abcdab -> length is 2 but the next char does not match. So, the pointer we were using to match from the front retraces by 1 char to check if it can match. i does not move forward yet.
        elif length > 0:
            length = lpsArr[length-1]
        else:
            lpsArr[i] = 0
            i += 1
    print(s)
    print(lpsArr)
    return lpsArr

def main():
    describe()

    input = "acbcabcbb"
    print("Input : " + input)
    print("Output : " + shortestPalindrome(input))
    print("\n\n")

    input = "bcaacb"
    print("Input : " + input)
    print("Output : " + shortestPalindrome(input))
    print("\n\n")

main()
