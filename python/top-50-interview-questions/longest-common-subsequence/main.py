def describe():
    desc = """
Problem : Given 2 strings str1 and str2, return the length of their longest common subsequence

Example :
    Input : "abdacbab", "acebfca"
    Output : 4 because longest common subsequence is "abca"

    Input : "cbebaff", "aeddbggf"
    Output : 3 because longest common subsequence is "ebf"

    Input : "abebafba", "cddghcd"
    Output : 0 because the longest common subsequence is 0
    """
    print(desc)

#  There is a string subsequence problem and we could possibly use that to get all subsequences of string 1
#  And then for each of those subsequences we can check if a subsequence from string 1 is a subsequence of string 2
#  and keep the one with the max length
#  If n is the length of string 1 and m is the length of string 2, we will get a time complexity of O((n+m) * 2^n)

#  The other more optimal/better solution is using 2 pointers and recursion
#  Take the strings "abdacbab", "acebfca" for example and 2 pointers ptr1 and ptr2
#    - First char matches - "a". lcs("abdacbab", "acebfca") = "a" + lcs("bdacbab", "cebfca")
#    - Second chars dont match, but the second char of one string can come later in the other string or vice versa. We dont know that.
#           So, we have 2 choices, move ptr1 forward once and get a solution and move ptr2 forward and get another solution
#           Now take the solution that gives you the subsequence with highest length
#           lcs("bdacbab", "cebfca") = max(lcs("dacbab", "cebfca"), lcs("bdacbab", "ebfca"))
def lcs(str1, str2):
    output = lcs_recursive(str1, str2)
    print(output)
    return len(output)

def lcs_recursive(str1, str2, output=""):
    if len(str1) == 0 or len(str2) == 0:
        return output

    if str1[0] == str2[0]:
        #  If the first chars of the strings match, then take that char into the solution
        #  and move the pointer forward for both the strings
        output += str1[0]
        output = lcs_recursive(str1[1:], str2[1:], output)
    else:
        #  Otherwise move the pointer of 1 string at a time forward and get both solutions.
        #  Then compare the results and take the one with the max length
        solution1 = lcs_recursive(str1[1:], str2, output)
        solution2 = lcs_recursive(str1, str2[1:], output)
        if len(solution1) >= len(solution2):
            output = solution1
        else:
            output = solution2

    return output



#  The recursion tree can contain n+m recursion calls. So, the time complexity is O(2^(n+m))
def more_concise_recursive_lcs(str1, str2, ptr1=0, ptr2=0):
    if ptr1 == len(str1) or ptr2 == len(str2):
        return 0
    elif str1[ptr1] == str2[ptr2]:
        return 1 + more_concise_recursive_lcs(str1, str2, ptr1+1, ptr2+1)
    else:
        return max(more_concise_recursive_lcs(str1, str2, ptr1+1, ptr2), more_concise_recursive_lcs(str1, str2, ptr1, ptr2+1))




def top_down_dp(str1, str2):
    #  Initialize a dp array such that in each cell it stores the max length of subsequence
    #  for indexes considered for str1 and str2 from poition i and j to the end
    #  ie, if i = 2, j = 3, then cell <i, j> will contain max subsequence length from postion i to end for string 1 and from position j to end for string 2
    #  Another way to think about this is that dp is storing the values of the sub problems of when iterating on positions <i,j> for string 1 and string 2
    #  This means that dp[0][0] will always contain the final result
    #  The initial values of the cells is 0
    dp = [[0 for j in range(len(str2))] for i in range((len(str1)))]
    top_down_dp_recursive_lcs(dp, str1, str2)
    return dp[0][0]

#  The recursion tree can contain n+m recursion calls. So, the time complexity is O(2^(n+m))
def top_down_dp_recursive_lcs(dp, str1, str2, i=0, j=0, maxlen = 0):
    if i == len(str1) or j == len(str2):
        return 0
    else:
        if dp[i][j] > 0:
            return dp[i][j]
        else:
            if str1[i] == str2[j]:
                val = 1 + top_down_dp_recursive_lcs(dp, str1, str2, i+1, j+1)
                dp[i][j] = val
                return val
            else:
                val = max(top_down_dp_recursive_lcs(dp, str1, str2, i+1, j), top_down_dp_recursive_lcs(dp, str1, str2, i, j+1))
                dp[i][j] = val
                return val

#  This solution can be improved further by memoization or dynamic programming
def main():
    describe()

    inputs = ["abdacbab", "acebfca"]
    print("Inputs : " + str(inputs))
    print("Outputs : " + str(top_down_dp(inputs[0], inputs[1])))
    print("\n\n")

    inputs = ["cbebaff", "aeddbggf"]
    print("Inputs : " + str(inputs))
    print("Outputs : " + str(top_down_dp(inputs[0], inputs[1])))
    print("\n\n")

    inputs = ["abebafba", "cddghcd"]
    print("Inputs : " + str(inputs))
    print("Outputs : " + str(top_down_dp(inputs[0], inputs[1])))
    print("\n\n")

main()
