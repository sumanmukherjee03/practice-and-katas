#  Problem : Given a string create a function that returns the first repeating character
#  If such a character does not exist return the null character '\0'

# This is the brute force solution.
#  Iterate over all elements of the string char by char
#  At any given point if a char exists in the substring just before that character then it is a repeated char
# The time complexity of this solution is O(n^2)
def firstRepeatingCharacter01(str):
    for i in range(len(str)):
        for j in range(i):
            if str[j] == str[i]:
                return str[i]
    return '\0'

#  Another idea could be to sort the characters by their ascii code. And simply check the previous character to see if there's any repeatation
#  But that does not work for our problem because we are looking for the *first* repeating char.
#  And sorting will break the order

# This is a better solution since it uses a hash table to store information about the occurance of a character.
#  The time complexity of this solution is O(n)
def firstRepeatingCharacter02(str):
    visited = {}
    for char in str:
        if visited.get(char):
            return char
        else:
            visited[char] = True
    return '\0'
