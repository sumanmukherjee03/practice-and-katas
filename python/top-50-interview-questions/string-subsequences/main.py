#  Problem : Given a string str write a function to return all possible subsequences, the order doesnt matter.
#  For ex: str = "abcd"
#  Output : ["abcd", "abc", "abd", "ab", "acd", "ac", "ad", "a", "bcd" .... "c", "d", ""]

# Solution - There are 2 possibilities with each character in the string. You can either take it or reject it.
# You stop your recursive process when you have considered all the characters.
#                                     ""
#                                    /  \
#                                  /      \
#                                "a"      ""             <- take or reject char "a"
#                                /\        /\
#                              /    \    /    \
#                           "ab"   "a"  "b"   ""        <- take or reject char "b"
#                           / \    / \   / \   / \
#                         /    \  /   \ /   \ /   \
#                       abc   ab ac   a bc  b c    ""   <- take or reject char "c"
#
# Time complexity is O(n * 2^n)
def getSubsequences(str):
    out = []
    def recursiveSubsequences(s, index = 0, subsequence = ""):
        if index == len(str):
            out.append(subsequence)
        else:
            recursiveSubsequences(s, index+1, subsequence+s[index])
            recursiveSubsequences(s, index+1, subsequence)
    recursiveSubsequences(str)
    return out
