def describe():
    desc = """
Problem : Given an integer n, return the number of strings of length n that consist of vowels and are lexicographically sorted.inserting a char

Example :
    Input : n = 2
    Output : 15
    Explanation : aa, ae, ai, ao, au, ee, ei, eo, eu, ii, io, iu, oo, ou, uu

------------------
    """
    print(desc)






def basic_solution_count_1(n):
    result = permutations(n)
    return len(result)

def permutations(n):
    vowels = ['a', 'e', 'i', 'o', 'u']
    if n == 0:
        return []
    elif n == 1:
        return vowels
    else:
        prev = permutations(n-1)
        result = []
        for i in range(0, len(vowels)):
            c = vowels[i]
            for s in prev:
                if s[0] >= c:
                    result.append(c+s)
        return result









#  It is important to note here that for
#    n=2
#       - if we pick a, there's 5 possibilities
#       - if we pick e, there's 4 possibilities
#       - if we pick i, there's 3 possibilities
#       - if we pick o, there's 2 possibilities
#       - if we pick u, there's 1 possibilities
#
#  Time complexity of this solution is O(5^n)
#
#  The recusive calls become like this
#       - func(3,a) + func(3,e) + ...... + func(3,u)
#       - func(3,a) expands to func(2,a) + func(2,e)+ ...... func(2,u)
#       - func(3,e) expands to func(2,e) + func(2,i)+ ...... func(2,u)
#       - func(3,i) expands to func(2,i) + func(2,o)+ ...... func(2,u)
#
#  Next level
#       - func(2,a) will expand to func(1,a)+func(1,e)+ .... + func(1,u). return 5 strings
#       - func(2,e) will return 4 strings
#
#  and so on...
def basic_solution_count_2(n):
    result = []
    for i in range(5):
        result += improved_permutations(n, i)
    return len(result)

def improved_permutations(n, currentIndex = 0):
    vowels = ['a', 'e', 'i', 'o', 'u']
    if n == 0 or currentIndex >= len(vowels):
        return []
    elif n == 1:
        return vowels[currentIndex]
    else:
        result = []
        prev = []
        for j in range(currentIndex, len(vowels)):
            prev += improved_permutations(n-1, j)

        c = vowels[currentIndex]
        for s in prev:
            result.append(c+s)

        return result









def count(n, lastChar=''):
    #  Base case is where no more chars to add are there so, we return 1
    if n == 0:
        return 1
    else:
        nb = 0
        vowels = ['a', 'e', 'i', 'o', 'u']
        for v in vowels:
            #  We keep track of the last character that was considered.
            #  That way we only consider the permutations for the vowels that are greater than equal to the last character
            if lastChar <= v:
                nb += count(n-1, v)
        return nb








def bottom_up_dp(n):
    #  Here each cell dp[i][j] represents the number of valid strings of size i that start with the vowel j
    #  ie strings of size dp[2][2] represents strings of size 2 that start with character e. NOT with character i.
    dp = [[0 for j in range(5+1)] for i in range(n+1)]

    #  number of strings of size 1 that start with various vowels
    for j in range(1, 5+1):
        dp[1][j] = 1

    for i in range(2, n+1):
        for j in range(1, 5+1):
            #  number of strings of size i that start with character j
            dp[i][j] = sum(dp[i-1][j:])

    for row in dp:
        print(row)

    result = 0
    for j in range(0, 5+1):
        result += dp[n][j]
    return result




def main():
    describe()

    input = 2
    print("Input : " + str(input))
    print("Output : " + str(bottom_up_dp(input)))
    print("\n\n")

    input = 3
    print("Input : " + str(input))
    print("Output : " + str(bottom_up_dp(input)))
    print("\n\n")

    input = 9
    print("Input : " + str(input))
    print("Output : " + str(bottom_up_dp(input)))
    print("\n\n")

main()
