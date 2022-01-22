def describe():
    desc = """
Problem : Given an array of unsorted numbers, find all unique triplets in it that add up to zero
Example :
    Input: [-3, 0, 1, 2, -1, 1, -2]
    Output: [-3, 1, 2], [-2, 0, 2], [-2, 1, 1], [-1, 0, 1]
    Explanation: There are four unique triplets whose sum is equal to zero.

---------
    """
    print(desc)

#  Brute force explanation:
#  3 loops
#  Outer loop which takes 1 number at a time.
#  Inner loop which takes another number at a time.
#  Another nested loop which takes a number other than those 2 numbers.

#################################### ONE SOLUTION #################################

#  An improvement over this is to sort the array first
#  Then at least the complexity of the algo will be O(nlog(n))
#  Now if you loop once over each number and inside the nested loop apply the solution for pair with sum problem inside the loop
#  Then the complexity will become O(nlog(n)) + O(n^2), which combined together is O(n^2)
#  One important thing to keep in mind here is that we are looking for unique triplets. So, we have to avoid duplicate numbers in each triplet.
#  And since the array is gonna be sorted that is easy to skip.
def triplet_sum_zero(arr):
    out = {} # maintain a hash to hold the results because we dont want the same triplets
    arr.sort()
    for i in range(0, len(arr)):
        num = arr[i]
        remainder = -num
        two_pair_res = two_pair_sum(arr, remainder, num)
        for j in range(0, len(two_pair_res)):
            solution = two_pair_res[j]
            solution.append(num)
            solution.sort()
            key = ','.join(map(str, solution))
            if key not in out:
                out[key] = solution
    return out.values()

def two_pair_sum(arr, target, avoid_number):
    out = []
    head = 0
    tail = len(arr) - 1
    while head < tail:
        n1 = arr[head]
        n2 = arr[tail]
        #  If value at head is the same as the one we are trying to avoid, then go to next iteration
        if n1 == avoid_number:
            head += 1
            continue
        #  If value at tail is the same as the one we are trying to avoid, then go to next iteration
        if n2 == avoid_number:
            tail -= 1
            continue
        if n1 + n2 == target:
            out.append([n1, n2])
            head += 1
            tail -= 1
        elif n1 + n2 < target:
            head += 1 # we need a pair of numbers with a bigger sum
        else:
            tail -= 1 # we need a pair of numbers with a smaller sum
    return out

#################################### ALTERNATE SOLUTION #################################

#  Time complexity is O(n^2)
def improved_triplet_sum_zero(arr):
    arr.sort()
    triplets = []
    for i in range(0, len(arr)):
        if i > 0 and arr[i] == arr[i-1]: # skip same element to avoid searching for duplicate triplets
            continue
        search_pair(arr, -arr[i], i+1, triplets) # search from i+1 th index onwards because we have already looked for triplets with earlier indexes
    return triplets

def search_pair(arr, target, left, triplets):
    right = len(arr) - 1
    while (left < right):
        current_sum = arr[left] + arr[right]
        if current_sum == target: # found a triplet
            triplets.append([-target, arr[left], arr[right]])
            left += 1
            right -= 1
            while left < right and arr[left] == arr[left-1]:
                left += 1 # skip the same element to avoid duplicates in the triplet
            while left < right and arr[right] == arr[right+1]:
                right -= 1 # skip the same element to avoid duplicates in the triplet
        elif current_sum < target:
            left += 1 # we need a pair of numbers with a bigger sum
        else:
            right -= 1 # we need a pair of numbers with a smaller sum


def main():
    describe()
    input = [-3, 0, 1, 2, -1, 1, -2]
    print("Input : ", input)
    res = improved_triplet_sum_zero(input)
    print("Output : ", res)

main()
