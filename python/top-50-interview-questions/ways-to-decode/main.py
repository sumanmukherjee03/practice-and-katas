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



def optimizedWaysToDecode(s, i = 0):
    n = len(s)
    if n == 0 or (i < n and s[i] == "0"):
        return 0
    elif i >= n-1:
        return 1
    elif 10 <= int(s[i] + s[i+1]) <= 26:
        return optimizedWaysToDecode(s, i+1) + optimizedWaysToDecode(s, i+2)
    else:
        return optimizedWaysToDecode(s, i+1)


def main():
    describe()

    input = '6324120129'
    print("Input : " + input)
    print("Output : " + str(waysToDecode(input)))
    print("\n\n")

    input = '12'
    print("Input : " + input)
    print("Output : " + str(waysToDecode(input)))
    print("\n\n")

main()
