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
#  Another loop which takes a number other than those 2 numbers.

#  An improvement over this is to sort the array first
#  Then at least the complexity of the algo will be O(nlog(n))
#  Now if you loop once and inside the nested loop apply the solution for pair with sum problem inside the loop
#  then the complexity will become O(nlog(n)) + O(n^2)
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
            head += 1
        else:
            tail -= 1
    return out

def main():
    describe()
    input = [-3, 0, 1, 2, -1, 1, -2]
    res = triplet_sum_zero(input)
    print("Input : ", input)
    print("Output : ", res)

main()
