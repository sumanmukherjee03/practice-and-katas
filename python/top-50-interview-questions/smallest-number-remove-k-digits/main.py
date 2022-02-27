def describe():
    desc = """
Problem : Given a positive integer k and string num that represents a positive integer, return a string that is the smallest num after removing k digits.
          Note that input and output both cant contain leading 0 unless it is the number 0 itself.

Example :
    Input : "825563", k=2
    Output : "2553"

    Input : "83866", k=3
    Output : "36"


--------------
    """
    print(desc)

def smallestAfterRemoving(num, k):
    digits = list(num)
    return recursiveSolution(digits, k, 0, digits[:])

def recursiveSolution(digits, k, currentIndex=0, currentSolution=[]):
    #  print("Called with : <" + str(k) + "," + str(currentIndex) + "," + str(currentSolution) + ">")
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

main()
