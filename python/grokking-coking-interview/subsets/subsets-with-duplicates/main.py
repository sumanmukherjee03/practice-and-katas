def describe():
    desc = """
Problem : Given a set of numbers that might contain duplicates, find all of its distinct subsets.
Example :
    Input: [1, 3, 3]
    Output: [], [1], [3], [1,3], [3,3], [1,3,3]

----------
    """
    print(desc)

#  Implementation one which follows the same approach as distinct subsets but with an added condition. Still is O(n * 2^n)
def brute_force_find_subsets(nums):
    subsets = []
    if len(nums) == 0:
        return subsets
    subsets.append([])
    for i in range(0, len(nums)):
        for j in range(0, len(subsets)):
            subset = list(subsets[j]) # Make a copy of the subset
            subset.append(nums[i])
            if subset not in subsets:
                subsets.append(subset)
    return subsets

def optimized_find_subsets(nums):
    subsets = []
    nums.sort()
    if len(nums) == 0:
        return subsets
    subsets.append([])
    prev_subsets = []
    for i in range(0, len(nums)):
        temp_subsets = [] # Keep a loop local variable to keep track of the subsets getting added in this iteration
        if i > 0 and nums[i] == nums[i-1]:
            #  When you find a duplicate number only append the new number to the previous set of subsets added
            #  We do this because the number before this has already added itself to empty subset and single number subsets etc.
            #  If we revisit the same subsets with this duplicate number we will be recreating the same set of subsets again plus some more by adding the dup number to the previous set of subsets
            #  So, to prevent the duplication of subsets we only consider the previous set of subsets
            #  Also, we need to sort the numbers first so that duplicate numbers are always consecutive. Otherwise this strategy doesnt work
            for j in range(0, len(prev_subsets)):
                subset = list(prev_subsets[j]) # Make a copy of the subset
                subset.append(nums[i])
                temp_subsets.append(subset)
            subsets = subsets + temp_subsets
        else:
            for j in range(0, len(subsets)):
                subset = list(subsets[j]) # Make a copy of the subset
                subset.append(nums[i])
                subsets.append(subset)
                temp_subsets.append(subset)
        prev_subsets = temp_subsets # Update previous set of subsets added with the current set of subsets added
    return subsets

def find_subsets(nums):
    # sort the numbers to handle duplicates
    list.sort(nums)
    subsets = []
    subsets.append([])
    startIndex, endIndex = 0, 0
    for i in range(len(nums)):
        startIndex = 0
        # if current and the previous elements are same, create new subsets only from the subsets
        # added in the previous step
        if i > 0 and nums[i] == nums[i - 1]:
            startIndex = endIndex + 1

        # endIndex maintains the history where the last set of subsets ended
        #  Since we use the value of endIndex to generate the value of startIndex before we update the endIndex when we encounter duplicate numbers
        #  and since we dont update the endIndex value after adding the new subsets in this iteration, that's why on the next iteration if duplicate number is encountered
        #  the startIndex will only give us the values from where the last set of subsets started
        endIndex = len(subsets) - 1

        for j in range(startIndex, endIndex+1):
            # create a new subset from the existing subset and add the current element to it
            set1 = list(subsets[j])
            set1.append(nums[i])
            subsets.append(set1)
    return subsets

#  Time complexity is still O(n * 2^n)
def main():
    describe()
    print("Here is the list of subsets: " + str(find_subsets([1, 3, 3])))
    print("Here is the list of subsets: " + str(find_subsets([1, 5, 3, 3])))
main()
