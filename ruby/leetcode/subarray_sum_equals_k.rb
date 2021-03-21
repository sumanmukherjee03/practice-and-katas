# Given an array of integers nums and an integer k, return the total number of continuous subarrays whose sum equals to k.

# Example 1:

# Input: nums = [1,1,1], k = 2
# Output: 2
# Example 2:

# Input: nums = [1,2,3], k = 3
# Output: 2


# Constraints:

# 1 <= nums.length <= 2 * 104
# -1000 <= nums[i] <= 1000
# -107 <= k <= 107


# The idea behind this approach is as follows: If the cumulative sum(represented by sum[i] for sum up to i th index)
# up to two indices is the same, the sum of the elements lying in between those indices is zero.
# Extending the same thought further, if the cumulative sum up to two indices, say i and j is at a difference of k
# i.e. if sum[i] - sum[j] = k, the sum of elements lying between indices i and j is k.

# Based on these thoughts, we make use of a hashmap mapmap which is used to store
# the cumulative sum up to all the indices possible along with the number of times the same sum occurs.
# We store the data in the form: (sum_i, no. of occurrences of sum_i).
# We traverse over the array nums and keep on finding the cumulative sum.
# Every time we encounter a new sum, we make a new entry in the hashmap corresponding to that sum.
# If the same sum occurs again, we increment the count corresponding to that sum in the hashmap.
# Further, for every sum encountered, we also determine the number of times the sum sum-k has occurred already,
# since it will determine the number of times a subarray with sum k has occurred up to the current index.
# We increment the count by the same amount.

# After the complete array has been traversed, the count gives the required result.

# @param {Integer[]} nums
# @param {Integer} k
# @return {Integer}
def subarray_sum(nums, k)
  sum = 0
  h = {}
  h[0] = 1
  count = 0

  # If there is a sub array of sum - k, then there is a contiguous array of sum k
  # including this current element
  i = 0
  while i < nums.length do
    sum += nums[i]
    count += h[sum - k] if h.has_key?(sum - k)
    h[sum] = h.fetch(sum, 0) + 1
    i += 1
  end

  return count
end


# This solution is based on cumulative sum
# If the cumulative sum up to an element itself is equal to k, then there's a subarray from 0 to that index with sum k
# Also iterate from that index + 1 to the end to see if there's any subarray whose sum is equal to k as well

# @param {Integer[]} nums
# @param {Integer} k
# @return {Integer}

def subarray_sum(nums, k)
  sums = []

  # Find cumulative sum up to each index
  nums.each_with_index do |num, index|
    sums[index] = index - 1 >= 0 ? sums[index-1] + num : num
  end

  count = 0
  # Iterate over all indexes
  for i in (0..(nums.length-1)) do
    # If cumulative sum at one index is k then from 0 to that index, the elements of the array are of sum k
    count += 1 if sums[i] == k

    # Go over the cumulative sums from index + 1 to the end to see if the difference in sums equals k
    # Then subarray from index+1 to that new index is also of sum k.
    if i+1 <= nums.length - 1
      for j in ((i+1)..(nums.length-1)) do
        count += 1 if (sums[j] - sums[i] == k)
      end
    end
  end

  return count
end
