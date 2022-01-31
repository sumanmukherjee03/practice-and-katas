from collections import deque

def describe():
    desc = """
Problem: Given a set of distinct numbers, find all it's permutations. Total number of permutations is n!

---------
    """
    print(desc)

#  This is similar to the subsets problem except that each subset must contain all the elements
#  Consider this process :
#      1. If the given set is empty then we have only an empty permutation set: []
#      2. Let's add the first element (1), the permutations will be: [1]
#      3. Let's add the second element (3), the permutations will be: [3,1], [1,3]
#      4. Let’s add the third element (5), the permutations will be: [5,3,1], [3,5,1], [3,1,5], [5,1,3], [1,5,3], [1,3,5]

#  This is the iterative solution and is a bit brute
#  Time complexity is O(n* n!)
def permutations(nums):
    result = [] # This container maintains the permutations of a subset of digits at any given time
    if len(nums) == 0:
        return result
    for i in range(0, len(nums)):
        new_digit = nums[i]
        if len(result) == 0: # Initiate result with a number of only 1 digit and that can have only 1 permutation
            result = [[new_digit]]
        else: #  Otherwise given a previous set of permutations of n digits add a new digit to it and generate the permutations of n+1 digits
            temp = []
            for i in range(0, len(result)):
                digits = result[i]
                for pos in range(0, len(digits)):
                    #  Only if the position of the digit is 0 should we add the new digit at the begining of the number
                    #  Otherwise add the new digit after the current digit in the iteration
                    if pos == 0:
                        temp.append([new_digit] + digits)
                    temp.append(digits[0:pos+1] + [new_digit] + digits[pos+1:len(digits)])
            result = temp
    return result


#  A more elegant solution than the previous one using a queue
#  Time complexity is O(n* n!)
def find_permutations(nums):
    numsLength = len(nums)
    result = []
    permutations = deque()
    permutations.append([])
    # currentNumber tracks 1 single digit at any given time
    for currentNumber in nums:
            # we will take all existing permutations and add the current number to create new permutations
            n = len(permutations)
            for _ in range(n):
                # pop a permutation from the queue one at a time and create a new permutation from it by adding the digit
                oldPermutation = permutations.popleft()
                # create a new permutation by adding the current number at every position
                for j in range(len(oldPermutation)+1):
                    newPermutation = list(oldPermutation) # copy the old permutation
                    newPermutation.insert(j, currentNumber) # insert new digit at jth position
                    if len(newPermutation) == numsLength: # If you have reached the desired number of digits in the new permutation, then add it to results
                        result.append(newPermutation)
                    else: # Otherwise push it into the queue
                        permutations.append(newPermutation)
    return result



#  A more elegant recursive solution that iterates on each permutation with a new digit
def generate_permutations(nums):
  result = []
  generate_permutations_recursive(nums, 0, [], result)
  return result

def generate_permutations_recursive(nums, index, currentPermutation, result):
    if index == len(nums):
        result.append(currentPermutation)
    else:
        # create a new permutation by adding the current number at every position
        for i in range(len(currentPermutation)+1):
            newPermutation = list(currentPermutation)
            newPermutation.insert(i, nums[index])
            generate_permutations_recursive(nums, index + 1, newPermutation, result) # with each permutation recursively add digits

def main():
    describe()
    input = [1, 3, 5]
    print("Input : " + str(input))
    print("Here are all the permutations: " + str(generate_permutations(input)))

main()
