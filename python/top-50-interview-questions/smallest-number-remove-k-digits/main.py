from collections import deque

def describe():
    desc = """
Problem : Given a positive integer k and string num that represents a positive integer, return a string that is the smallest num after removing k digits.
          Note that input and output both cant contain leading 0 unless it is the number 0 itself.

Example :
    Input : "825563", k=2
    Output : "2553"

    Input : "83866", k=3
    Output : "36"

    Input : "4000235", k=1
    Output : "235"


--------------
    """
    print(desc)

def smallestAfterRemovingSolution1(num, k):
    digits = list(num)
    return recursiveSolution(digits, k, 0, digits[:])

#  As we process each digit, we have 2 choices, either keep it in the solution or reject it.
#  Maintain the current solution as a list of chars.
#  If we reject a digit in the current solution mark it as 'x'.
#  That way at the end when we have performed 'k' removals, from the list of chars filter out the chars which are 'x' and get the final string which will be the solution
def recursiveSolution(digits, k, currentIndex=0, currentSolution=[]):
    #  print("Recursion called with : <" + str(k) + "," + str(currentIndex) + "," + str(currentSolution) + ">")
    if len(digits) == 0:
        return ''
    if k == 0 or currentIndex == len(digits)-1:
        temp = [i for i in currentSolution if i != 'x']
        return ''.join(temp)
    else:
        solution1 = recursiveSolution(digits, k-1, currentIndex+1, currentSolution[0:currentIndex] + ['x'] + currentSolution[currentIndex+1:])
        solution2 = recursiveSolution(digits, k, currentIndex+1, currentSolution)
        if len(solution1) == 0:
            solution1 = '0'
        if len(solution2) == 0:
            solution2 = '0'
        solution = min(int(solution1), int(solution2))
        return str(solution)







#  Some properties to note when forming a number from digits.
#  Lets say we have some random digits 3, 6, 5, 9.
#  The smallest number will be the number where these digits are sorted, ie, 3569. Although there are n! combinations of it mathematically.
#  It is also necessary to realize that the position of a digit determines it's weight.
#  Meaning even if 9 is greater than 3, it has more weight because 3 is left of 9. so, the more left a digit is, the more weight it has.
#
#  Armed with this knowledge, when we attack the problem above, there is a greedy approach we can take.
#  As we traverse the digits, if the current digit we are considering is smaller than the digit that came just before it,
#  we can delete that digit and replace it with the current digit. Repeatedly do this until either we have done k deletions
#  or have reached a previous digit that is smaller than the current digit.
#  For example : Take 26378491 and k = 3
#  Keep 2, keep 6 - 26
#  When considering 3, 3 is smaller than 6. So, remove 6 and the current number becomes 23. k becomes 2.
#  Next take 7, take 8. Current number is 2378.
#  When we consider 4, again 4 is smaller than 8. So, remove 8. 4 is smaller than 7, so remove 7.
#  k becomes 0 and also 3 < 4. So, we stop. Number is 234. k = 0.
#  Now, take the rest of the digits as is - 23491. That's our final solution.
#
#  As you can see the greedy solution is based on the fact that we are aiming to reach a sorted digits solution. Or at least part way there.
#  And we traverse the digits left to right because the left most digits have the most weight.
#
#  Since we are inspecting the last element and then the one before it, we need a last in first out data structure. Stack is perfect here.
#  Time complexity of this is O(n)
def smallestAfterRemoving(num, k):
    if k == len(num):
        return '0'
    digits = list(num)
    stack = deque()
    for d in digits:
        while len(stack) > 0 and k > 0 and int(stack[-1]) > int(d):
            stack.pop()
            k -= 1
        stack.append(d)

    #  This is to handle the case where the input string is sorted or the elements in the stack at the surrent state are sorted.
    #  For example you could end up having 2 same digits. They might end up in the stack. Which is fine, but k might not be 0 yet.
    #  So, keep popping until k gets to 0
    while k > 0 and len(stack) > 0:
        stack.pop()
        k -= 1

    #  Due to this process after removing k digits, you might endup having leading 0's in front of the stack.
    #  And that we do not want. So, we should keep popping from the front, ie left until all leading 0's are removed.
    #  NOTE : Unnecessary tip. To reverse an array in python arr[::-1] . Not that we are using that tip here.
    while len(stack) > 0 and int(stack[0]) == 0:
        stack.popleft()

    if len(stack) == 0:
        return '0'
    else:
        return "".join(stack)



def main():
    describe()

    input = "825563"
    k = 2
    print("Input : " + input + ", " + str(k))
    print("Output : " + smallestAfterRemoving(input, k))
    print("\n\n")

    input = "83866"
    k = 3
    print("Input : " + input + ", " + str(k))
    print("Output : " + smallestAfterRemoving(input, k))
    print("\n\n")

    input = "4000235"
    k = 1
    print("Input : " + input + ", " + str(k))
    print("Output : " + smallestAfterRemoving(input, k))
    print("\n\n")

main()
