#  Problem : Given a string with brackets [{()}], create a function that checks
#  if string is valid, ie all brackets are closed.

# If it's an opening bracket push to a stack
# If it's a closing bracket, pop from the stack and match brackets
# If there's a closing bracket but there is no opening bracket, then the string is not valid.
# If there are opening brackets left in the stack after iterating on the string it means there
# were more opening brackets than closing brackets, so string is not valid
# Time complexity is O(n)
def isValid01(str):
    stack = []
    for c in str:
        if c in ["(", "{", "["]:
            stack.append(c)
        else:
            if len(stack) == 0:
                return False
            else:
                prev = stack.pop()

            if c == ")" and prev != "(":
                return False
            elif c == "}" and prev != "{":
                return False
            elif c == "]" and prev != "[":
                return False

    return len(stack) == 0


# Same solution as above but with a slightly different implementation
def isValid02(str):
    bracketsMap = {"(": ")", "{": "}", "[": "]"}
    openingBrackets = ["(", "{", "["]
    stack = []
    for c in str:
        if c in openingBrackets:
            stack.append(c)
        elif len(stack) > 0 and bracketsMap[stack[-1]] == c:
            stack.pop()
        else:
            return False
    return len(stack) == 0
