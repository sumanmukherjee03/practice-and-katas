def describe():
    desc = """
Problem : Given 2 words as strings word1 and word2, return the number of operations required to convert word1 to word2
    This value is also called the Levenshtein distance.

Example :
    Input : str1 = "inside", str2 = "index"
    Output : 3
    Explanation : "inside" -> "index" requires 3 operations, remove 's', remove 'i' and insert 'x'
        (inside -> inide -> inde -> index)

    Input : str1 = "eagles", str2 = "algiers"
    Output : 4
    Explanation : "eagles" -> "algiers" requires 4 operations, remove 'e', insert 'l', replace second 'l' by 'i', insert 'r'
        (eagles -> agles -> algles -> algies -> algiers)

------------------------
    """
    print(desc)


#  Time complexity is O(3^(n+m))
def min_distance(str1, str2, i = 0, j = 0):
    #  If you have reached the end of both the first and the second string then the operation is complete
    if i == len(str1) and j == len(str2):
        return 0

    #  If you have reached the end of the second string but not the first that means you just have to remove the extra chars from str1
    elif j == len(str2) and i < len(str1):
        return len(str1) - i

    #  If you have reached the end of the first string but not the second that means you just have to insert the extra chars from str2 to end
    elif i == len(str1) and j < len(str2):
        return len(str2) - j

    else:
        #  Check the char at the current index for both strings at the moment.
        #  If both chars are the same, move both the pointers forward by 1.
        #  But no operation is performed, so we dont add 1 to the return statement.
        if str1[i] == str2[j]:
            return min_distance(str1, str2, i+1, j+1)
        else:
            #  If the chars at the current pointers did not match then there are 3 possible choices.
            #       - remove a char from str1
            #       - insert a char right before the current char in str1
            #       - replace the current char of str1
            #  We dont know which one of these operations will yield a result. So, we try all 3 and take the one with the min count
            sol1 = min_distance(str1[0:i] + str1[i+1:], str2, i, j) # remove char at i from str1
            sol2 = min_distance(str1[0:i] + str2[j] + str1[i:], str2, i+1, j+1) # insert char at j from str2 before char at i in str1
            sol3 = min_distance(str1[0:i] + str2[j] + str1[i+1:], str2, i+1, j+1) # replace char at position i of str1 with char at position j of str2
            return 1 + min(sol1, sol2, sol3)



#  Time complexity is O(3^(n+m))
def min_distance_concise(str1, str2, i=0, j=0):
    if i == len(str1):
        return len(str2)-j
    elif j == len(str2):
        return len(str1)-i
    elif str1[i] == str2[j]:
        return min_distance_concise(str1, str2, i+1, j+1)
    else:
        #  If the chars at the current pointers did not match then there are 3 possible choices.
        #       - imagine removing a char from str1 and move the pointer i forward
        #       - imagine inserting a char right before the current char in str1. That means we still want to process char at i of str1. So, dont move i, but move forward by 1.
        #       - imagine replacing the current char of str1. Means move both i and j by 1.
        #  We dont know which one of these operations will yield a result. So, we try all 3 and take the one with the min count.
        #  Either ways whichever operation is chosen we add 1 to that min distance and return that.
        return 1 + min(min_distance_concise(str1, str2, i+1, j),
                min_distance_concise(str1, str2, i, j+1),
                min_distance_concise(str1, str2, i+1, j+1))



#  Time complexity of this solution is O(n*m)
def bottom_up_dp(str1, str2):
    if len(str1) == 0 and len(str2) > 0:
        return len(str2)
    if len(str1) > 0 and len(str2) == 0:
        return len(str1)

    n, m = len(str1), len(str2)
    #  In this dp array each cell dp[i][j] represents the number of operations required to convert str1[0:i+1] into str2[0:j+1]
    dp = [[0 for j in range(m+1)] for i in range(n+1)]

    #  Find the number of operations required for various lengths of str1 to be converted to an empty string representing str2
    for i in range(1,n+1):
        dp[i][0] = i

    #  Find the number of operations required for an empty string representing str1 to be converted to various lengths of str2
    for j in range(1,m+1):
        dp[0][j] = j


    for i in range(1, n+1):
        for j in range(1, m+1):
            #  If the char at the last index of each string matched, then we must have moved both pointers i and j by 1
            if str1[i-1] == str2[j-1]:
                dp[i][j] = dp[i-1][j-1]
            else:
                #  Otherwise, we take the min of of the last choice plus 1
                #  And the last choice can be one of these :
                #  - removing a char from str1 meaning i moved but j remained the same
                #  - inserting a char into str1 meaning j moved but we keep processing from the same index of i
                #  - replacing a char in str1 meaning both i and j had moved by 1
                dp[i][j] = min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])+1

    for row in dp:
        print(row)

    return dp[n][m]



def main():
    describe()

    input = ["inside", "index"]
    print("Input : " + str(input))
    print("Output : " + str(bottom_up_dp(input[0], input[1])))
    print("\n")

    input = ["eagles", "algiers"]
    print("Input : " + str(input))
    print("Output : " + str(bottom_up_dp(input[0], input[1])))
    print("\n")

main()
