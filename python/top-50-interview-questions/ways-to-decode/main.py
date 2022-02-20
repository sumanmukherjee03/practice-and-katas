import string

def describe():
    desc = """
Problem : Given a string made of digits, create a function that returns how many ways we can decode it.
          It is given that 1 decodes to A, 2 decodes to B and so on upto 26.

Example :
        Input : 6324120129
        Output : We can decode this as FCBDATABI (6 3 2 4 1 20 1 2 9) or FCBDATLI (6 3 2 4 1 20 12 9)
        or FCXATABI (6 3 24 1 20 1 2 9) or FCXATLI (6 3 24 1 20 12 9) . So output is 4

        Input : 12
        Output : We can decode this as AB (1 2) or L (12)

--------------------
    """
    print(desc)

def waysToDecode(s):
    if len(s) == 0:
        return 0

    chars = list(s)
    digits = []
    for i in range(0, len(chars)):
        d = int(chars[i])
        if i < len(chars)-1 and d > 0 and int(chars[i+1]) == 0:
            digits.append(d*10)
        elif d > 0:
            digits.append(d)

    result = []
    recursiveWaysToDecode(result, digits)
    print("Possible decoded strings : " + str(result))
    return len(result)

#  Pass around result as a reference and keep appending the final decoded strings to it.
#  currentIndex holds the current index being processed.
#  currentResult holds the state of the current decoded string.
#  Time complexity of this recursive solution :
#    T(0) = 1
#    T(n) = T(n-1) + T(n-2)
#    which makes it a combinatorial series. So final T(n) = O(phi ^ n) where phi is the golden ratio
def recursiveWaysToDecode(result, digits, currentIndex = 0, currentResult = ""):
    letters = [""] + list(string.ascii_uppercase)

    if currentIndex >= len(digits):
        result.append(currentResult)
        return

    #  At each index being processed there's 2 choices. Take that digit stand alone.
    #  OR, combine it with the next digit and if it is less than 26 then it can become a different alphabet.
    recursiveWaysToDecode(result, digits, currentIndex+1, currentResult + letters[digits[currentIndex]])
    if currentIndex+1 <= len(digits)-1:
        currentDigit = digits[currentIndex]
        nextDigit = digits[currentIndex+1]
        num = currentDigit * 10 + nextDigit
        if num <= 26:
            recursiveWaysToDecode(result, digits, currentIndex+2, currentResult + letters[num])




#  This is the same recusive solution as before, just slightly smaller and more optimized code.
#  This is also of time complexity O(phi ^ n)
def betterWaysToDecode(s, i = 0):
    n = len(s)
    if n == 0 or (i < n and s[i] == "0"):
        return 0
    elif i >= n-1:
        return 1
    elif 10 <= int(s[i] + s[i+1]) <= 26:
        return optimizedWaysToDecode(s, i+1) + optimizedWaysToDecode(s, i+2)
    else:
        return optimizedWaysToDecode(s, i+1)




#  by applying dynamic programming to the slution above we can reduce the recalculation of the subproblems that we are solving recursively
#  Time complexity is O(n)
def bottomUpDP(s):
    digits = [int(c) for c in list(s)]

    n = len(digits)
    if n == 0 or digits[0] == 0:
        return 0

    #  At each index we are going to consider the number of ways to decode the string upto that current index only.

    #  We initialize beforePrevious to 1 instead of 0 like the fibinacci solution below because consider the case of 1 2.
    #  When we are considering the number 2 at index 1, the previous would have been 1 but if beforePrevious was 0 then the number of ways to decode would have been 1 and not 2 for the number at index 1
    #  That is why we initialize the beforePrevious with 1
    beforePrevious = 1
    previous = 1 # For only 1 char, the number of ways to decode the string is only 1
    for i in range(1, n):
        # Lets say at this current index the number of ways to decode the string is 0.
        #  Now there are 2 possibilities
        #  -  One is we consider just adding this digit.
        #  -    In this case, we get an additional number of ways to decode the same as the number of ways to decode up until the last character
        #  -  Other is if we join with the previous digit and get a number less than 26.
        #  -    In this case, we get an additional number of ways to decode the same as the number of ways to decode up until the before last character
        #  Say for a string 6 3 2 4
        #    The ways to decode up until each char would be 1 1 1 2
        current = 0
        if digits[i] > 0:
            current += previous
        if 10 <= digits[i-1]*10 + digits[i] <= 26:
            current += beforePrevious
        beforePrevious = previous
        previous = current
    return previous


#  The solution above follows the similar pattern as solving fibinacci with dynamic programming
def fibonacci(n):
    dp = [0] * (n+1)
    dp[0] = 0
    dp[1] = 1
    for i in range(2, n+1):
        dp[i] = dp[i-1] + dp[i-2]
    return dp[n]

#  Same fibonacci numbers but without the use of array
def fibonacciNoArray(n):
    beforePrevious = 0
    previous = 1
    for i in range(2, n+1):
        current = beforePrevious + previous
        beforePrevious = previous
        previous = current
    return previous



def main():
    describe()

    input = '6324120129'
    print("Input : " + input)
    print("Output : " + str(bottomUpDP(input)))
    print("\n\n")

    input = '12'
    print("Input : " + input)
    print("Output : " + str(bottomUpDP(input)))
    print("\n\n")

main()
