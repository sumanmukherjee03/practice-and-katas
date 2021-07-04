#  Problem : Given a string with alphabetical characters only create a function that
#  returns the length of the longest substring without repeating characters.
#  Ex : "abcdbeghef"
#  Output : 6 -> because the longest substring without repeating characters is "cdbegh"


#  The brute force solution is to find all possible substrings with no repeating characters
#  and return the length of the substring whose length is maximum
#  However the time complexity of this is O(n^3)
def longestSubstringWithoutRepeating01(str):
    if len(str) <= 1:
        return len(str)

    def hasNoRepeatingChars(str):
        charCount = [0] * 128 # Initialize an array for all the ascii chars upto code 128 with char count set to 0
        for char in str:
            charCount[ord(char)] = charCount[ord(char)] + 1
            if charCount[ord(char)] > 1:
                return False
        return True

    maxLen = 0
    #  Keep 2 loops - one for tracking the start of the index from where we want to take the substrings
    #  and in a nested loop for each starting index get all possible substrings till the end of the string
    for i in range(len(str)):
        for j in range(i,len(str)):
            substr = str[i:j+1] # Get the substring from index i upto index j
            if hasNoRepeatingChars(substr):
                maxLen = max(maxLen, len(substr))
    return maxLen



def longestSubstringWithoutRepeating02(str):
    if len(str) <= 1:
        return len(str)

    dp = [] # Stores the local maximum substring upto index i
    dp.append(str[0])
    maxSubstr = dp[0]

    for i in range(1, len(str)):
        #  If you found a char that is already present in the previous longest substring without repeating chars
        #  you want to find the part of the substring that is after that repeated character for the new sequence of chars
        if str[i] in dp[i-1]:
            strSinceLastOccuranceOfChar = dp[i-1].split(str[i])[-1]
            dp.append(strSinceLastOccuranceOfChar + str[i])
        else:
            dp.append(dp[i-1] + str[i])

        if len(dp[i]) > len(maxSubstr):
            maxSubstr = dp[i]

    return len(maxSubstr)



#  We maintain 2 pointers - one for start of the substring and one for the end of the substring
#  We maintain an array of ascii codes of chars where we store the index of the last occurance of a character.
#  The start pointer of the substrings begin at 0 and the stop pointer keeps moving forward
#  as long as there is no repeating char which we can find from the array maintaining the last seen position of char.
#  If a repeating char is encountered, it means the sliding window needs to move, as in the start pointer
#  needs to move after the repeated char.
#  Time complexity of this traversal technique is O(n)
def longestSubstringWithoutRepeating03(str):
    if len(str) <= 1:
        return len(str)

    charMap = [-1] * 128 # For each character store the index of the latest occurance of that character
    start = 0
    stop = 0
    maxLen = 0

    while stop < len(str):
        # If the char at stop pointer exists in the charMap, ie, it is not -1
        #  that means there is a repeated char. And the new start has to be moved to the position after the
        #  occurance of that character.
        if charMap[ord(str[stop])] >= start:
            start = charMap[ord(str[stop])] + 1
        charMap[ord(str[stop])] = stop
        substr = str[start:stop+1]
        maxLen = max(len(substr),maxLen)
        stop += 1

    return maxLen



#  This is exactly like the previous solution except that it is using a for loop for moving the stop pointer ahead
def longestSubstringWithoutRepeating04(str):
    if len(str) <= 1:
        return len(str)

    charMap = [-1] * 128 # For each character store the index of the latest occurance of that character
    start = 0
    maxLen = 0

    for i in range(len(str)):
        # If the char at i exists in the charMap, ie, it is not -1
        #  that means there is a repeated char. And the new start has to be moved to the position after the
        #  occurance of that character.
        if charMap[ord(str[i])] >= start:
            start = charMap[ord(str[i])] + 1
        charMap[ord(str[i])] = i
        maxLen = max(maxLen, i+1-start)

    return maxLen
