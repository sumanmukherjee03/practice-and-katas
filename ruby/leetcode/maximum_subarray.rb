# Given an integer array nums, find the contiguous subarray (containing at least one number) which has the largest sum and return its sum.

# Example 1:

# Input: nums = [-2,1,-3,4,-1,2,1,-5,4]
# Output: 6
# Explanation: [4,-1,2,1] has the largest sum = 6.
# Example 2:

# Input: nums = [1]
# Output: 1
# Example 3:

# Input: nums = [5,4,-1,7,8]
# Output: 23


# Constraints:

# 1 <= nums.length <= 3 * 104
# -105 <= nums[i] <= 105

# @param {Integer[]} nums
# @return {Integer}
def max_sub_array(nums)
  i = 0

  # This keeps track of the local maximum at an index
  # This can be the local maximum till the previous index + the current element if that is bigger than the current element itself
  # Or it is the current element.
  # This accounts for negative numbers.
  max_till_current = 0

  max = nil # This tracks the global maximum found so far
  while i < nums.length do
    if max_till_current + nums[i] > nums[i]
      max_till_current = max_till_current + nums[i]
    else
      max_till_current = nums[i]
    end
    nums[i] = max_till_current
    if !max || max_till_current > max
      max = max_till_current
    end
    i += 1
  end
  return max
end

def max_sub_array_another_form_of_dynamic_programming(nums)
  dp = Array.new(nums.length){Array.new(nums.length, 0)}
  max = {begin: 0, end: 0}

  # Think about you are iterating over the rows of the dp array
  j = 0
  while j < nums.length do
    # Think about you are iterating over the columns of the dp array upto the diagonal
    i = 0
    while i <= j do
      if j == i
        dp[i][j] = nums[i]
      else
        dp[i][j] = dp[i][j-1] + nums[j]
      end

      if dp[max[:begin]][max[:end]] < dp[i][j]
        max[:begin] = i
        max[:end] = j
      end

      i += 1
    end
    j += 1
  end

  res = nums[max[:begin]..max[:end]]
  return res.inject(0, :+)
end
