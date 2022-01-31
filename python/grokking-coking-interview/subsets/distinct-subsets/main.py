def describe():
    desc = """
Problem : Given a set of elements find all it's distinct subsets.

Example:
    Input : [1,3]
    Output : [], [1], [3], [1,3]

----------------
    """
    print(desc)

#  Follow the BFS approach
#  start with an empty set.
#  Iterate over the given list of numbers and add the number at current index to all the elements of the existing subsets.
#
#  In each set the number of elements double. So we have a total of 2^n subsets, ie, representing all the permutations.
#  Time complexity is O(n* 2^n)
def find_subsets(nums):
    subsets = []
    if len(nums) == 0:
        return subsets
    subsets.append([])
    for i in range(0, len(nums)):
        for j in range(0, len(subsets)):
            #  Remember that python references list by pointers
            #  So do not append to an existing list. But, rather copy it to a new list and append it.
            subset = subsets[j].copy() #  This also creates a new list - list(subsets[i])
            subset.append(nums[i])
            subsets.append(subset)
    return subsets

def main():
    describe()
    print("Here is the list of subsets: " + str(find_subsets([1, 3])))
    print("Here is the list of subsets: " + str(find_subsets([1, 5, 3])))

main()
