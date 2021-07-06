#  Problem : Given 2 strings str1 and str2 find the first index where we can find str2 in str1

# This solution will work for the most part. It's slightly high in time complexity but something we can work with.
#  Time complexity is O(nm)
def substrIndex01(str1, str2):
    if len(str1) < len(str2):
        return -1
    if len(str1) == len(str2):
        return 0 if str2 == str1 else return -1
    if len(str2) = 0:
        return 0
    len1 = len(str1)
    len2 = len(str2)
    found = -1
    # Iterate on the bigger string until there's at least enough characters left to match the smaller string
    for i in range(len1-len2+1):
        found = True
        #  From the ith position of first string match char by char with 2nd string
        #  If a mismatch is found reject and come out of the loop and start with the next ith character again.
        #  Think of it like a rubber band with the start end fixed. It increases in length as long as there are matches.
        #  In case of a mismatch it starts over again from the next char in the bigger string, ie the fixed start end moves 1 char.
        for j in range(len2):
            if str1[i+j] != str2[j]:
                found = False
                break
        if found:
            return i
    return -1




#  This solution is an improvement over the naive algorithm mentioned above.
#  It is a direct implementation of the KMP (Knuth-Morris-Pratt) algorithm.
#  Calculate the longest prefix that is also a suffix in a pattern.
#  If we can calculate the LPS in a substring or pattern then we dont need to rewind the comparison
#  from 0th position of the pattern and i+1 th position of the string.
#  ie, if a substring has occured more than once in the pattern to be matched, no need to move i back all the time.
#
#  LPS table example :
    #  a b c d a b e a b c
    #  0 0 0 0 1 2 0 1 2 3

    #  a a a a b a a c d
    #  0 1 2 3 0 1 2 0 0
# Time complexity of this algo is O(n + m) or else O(n)
def substrIndex02(str1, str2):
    def getLpsArr(str):
        lpsArr = [0] * len(str)
        length = 0
        i = 1
        while i < len(str):
            # This is the simple case. As long as chars are matching
            # LPS[i] stores the number of chars right before it that have matched chars from 0th index at the ith index
            # So, if LPS[5] is 2, it means starting from index 4 of the string 2 chars match the start
            # ie str[0] matches str[4] and str[1] matches str[5]
            if str[i] == str[length]:
                length += 1
                lpsArr[i] = length
                i += 1
            elif length > 0:
                length = lpsArr[length-1]
            # This is the default case for filling up theLPS array. It is always 0 by default
            else:
                lpsArr[i] = 0
                i += 1
        return lpsArr
    pass

    # Once you get the LPS array for the pattern that you are trying to match with the given string
    # maintain 2 pointers i and j for searching. i for the target string and j for the pattern.
    # If there are matches move i and j forward. When there is a character mismatch find the value of LPS[j-1]
    # and take j back to that index in the pattern. Dont move i forward yet. Start matching again with i and j.
    # If j has become 0, ie come to the start of the pattern and there is still no match with what's at ith index, move i forward.
    n = len(str1)
    m = len(str2)
    if m > n:
        return -1
    if m == n:
        return 0 if str2 == str1 else return -1
    if str2 = "":
        return 0
    lpsArr = getLpsArr(str2)
    i = 0
    j = 0
    while i < n and j < m:
        if str[i] == str[j]:
            i += 1
            j += 1
        elif j > 0:
            j = lpsArr[j-1]
        else:
            i += 1
    return -1 if j < m else i - j
