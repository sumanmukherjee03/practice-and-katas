# Given an integer array nums, return an array answer such that answer[i] is equal to the product of all the elements of nums except nums[i].

# The product of any prefix or suffix of nums is guaranteed to fit in a 32-bit integer.

# Example 1:

# Input: nums = [1,2,3,4]
# Output: [24,12,8,6]
# Example 2:

# Input: nums = [-1,1,0,-3,3]
# Output: [0,0,9,0,0]


# Constraints:

# 2 <= nums.length <= 105
# -30 <= nums[i] <= 30
# The product of any prefix or suffix of nums is guaranteed to fit in a 32-bit integer.


# Follow up:

# Could you solve it in O(n) time complexity and without using division?
# Could you solve it with O(1) constant space complexity? (The output array does not count as extra space for space complexity analysis.)#


# @param {Integer[]} nums
# @return {Integer[]}
def product_except_self(nums)
  # Instantiate with a base value of 1 to ensure that products_before[0] = 1
  # and product_after[nums.length - 1] = 1
  product_before = Array.new(nums.length, 1) # Instantiate with a value of 1 for the length of the array
  product_after = Array.new(nums.length, 1) # Instantiate with a value of 1 for the length of the array

  # Precompute product of elements on left of self
  i = 1
  while i < nums.length do
    product_before[i] = product_before[i-1] * nums[i-1]
    i += 1
  end

  # Precompute product of elements on right of self
  i = nums.length - 2
  while i >= 0 do
    product_after[i] = product_after[i+1] * nums[i+1]
    i -= 1
  end

  i = 0
  res = []
  while i < nums.length do
    res[i] = product_before[i] * product_after[i]
    i += 1
  end

  return res
end
