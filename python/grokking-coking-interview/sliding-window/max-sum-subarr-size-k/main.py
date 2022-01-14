def describe():
    str = """
Problem : Find the max sum of contiguous subarrays of size k
----------
    Ex : Input: [2, 1, 5, 1, 3, 2], k=3
    Output: 9
    Explanation: Subarray with maximum sum is [5, 1, 3]

#########
    """
    print(str)

#  Time complexity = O(n^2)
def find_max_sum_of_subarrays_brute_force(arr, k):
  max_sum = 0
  local_sum = 0
  for i in range(len(arr) - k + 1):
    local_sum = 0
    for j in range(i, i+k):
      local_sum += arr[j]
    max_sum = max(max_sum, local_sum)
  return max_sum

#  Time complexity = O(n)
def find_max_sum_of_subarrays(arr, k):
    maxSum = 0
    localSum = 0
    head = 0
    tail = 0
    for tail in range(len(arr)):
        elm = arr[tail]
        localSum += elm
        if tail - head >= k - 1:
            maxSum = max(localSum, maxSum)
            localSum -= arr[head]
            head += 1
    return maxSum

def main():
    describe()
    res = find_max_sum_of_subarrays([2, 1, 5, 1, 3, 2], 3)
    print(res)

main()
