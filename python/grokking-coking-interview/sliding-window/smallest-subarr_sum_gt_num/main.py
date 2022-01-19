import math

def describe():
    str = """
Problem : Find the smallest subarray of a given array whose sum is greater than a given number.

NOTE : In this sliding window example the window size is not fixed but based on a contraint
----------
    Ex : Input: [2, 1, 5, 2, 3, 2], S=7
    Output: 2
    Explanation: The smallest subarray with a sum greater than or equal to '7' is [5, 2]

#########
    """
    print(str)


#  The sliding window moves like a caterpillar
#  If local sum is less than k, move right side forward
#  If local sum is more than or equal to k move left side forward and keep doing that until the sum goes down below k
#  Although this looks like 2 loops the inner loop only goes through a very limited number of iterations. Total complexity is O(n)
def find_smallest_subarr_sum_gt_num(arr, k):
    window_sum = 0
    min_length = math.inf # NOTE : This is the way to represent infinity in python
    window_start = 0
    for window_end in range(0, len(arr)):
        window_sum += arr[window_end]  # add the next element
        # shrink the window as small as possible until the 'window_sum' is smaller than 's', all the time while moving the start pointer forward by 1
        # This movement is similar to that of a caterpillar.
        while window_sum >= k:
            min_length = min(min_length, window_end - window_start + 1) # NOTE : To find subarray length, always add 1 to the diff of indexes
            window_sum -= arr[window_start]
            window_start += 1
    # If at the end, no such subarray was found, return 0
    if min_length == math.inf:
        return 0
    return min_length


#  The sliding window moves like a caterpillar
#  If local sum is less than k, move right side forward
#  If local sum is more than or equal to k move left side forward
def alternate_find_smallest_subarr_sum_gt_num(arr, k):
    print("Input : {}, {}".format(arr, k))
    res = []
    head = 0 # This is where the subarray index starts and is always inclusive
    tail = 0 # This is where the subarray index ends and is always exclusive
    localSum = 0
    # Iterate as long as head does not cross tail
    #  And head inclusive can be the last index in the given array
    #  And tail exclusive can be of index 1 more than the last index of the given array
    while head <= tail and head <= len(arr)-1 and tail <= len(arr):
        if localSum < k:
            elm = arr[tail]
            localSum += elm
            tail += 1 # This is what makes the tail element to not be included in the subarray for local sum
        elif localSum >= k:
            if (tail - head) <= len(res):
                res = arr[head:tail]
            elif len(res) == 0:
                res = arr[head:tail]
            localSum -= arr[head]
            head += 1
    return res

def main():
    describe()
    res = find_smallest_subarr_sum_gt_num([2, 1, 5, 2, 8], 7)
    #  res = find_smallest_subarr_sum_gt_num([2, 1, 5, 2, 3, 2], 7)
    print(res)

main()
