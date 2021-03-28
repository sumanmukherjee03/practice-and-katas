# You are given an array of integers nums, there is a sliding window of size k which is moving from the very left of the array to the very right.
# You can only see the k numbers in the window. Each time the sliding window moves right by one position.

# Return the max sliding window.

# Example 1:
# Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
# Output: [3,3,5,5,6,7]
# Explanation:
# Window position                Max
# ---------------               -----
# [1  3  -1] -3  5  3  6  7       3
 # 1 [3  -1  -3] 5  3  6  7       3
 # 1  3 [-1  -3  5] 3  6  7       5
 # 1  3  -1 [-3  5  3] 6  7       5
 # 1  3  -1  -3 [5  3  6] 7       6
 # 1  3  -1  -3  5 [3  6  7]      7

# Example 2:

# Input: nums = [9,11], k = 2
# Output: [11]


#################### Caterpillar two pointer approach #####################

# @param {Integer[]} nums
# @param {Integer} k
# @return {Integer[]}
def max_sliding_window(nums, k)
  return [] if nums.length == 0
  return [] if k == 0
  return nums if nums.length == 1

  res = []
  left_ptr = 0
  right_ptr = (left_ptr + k) - 1
  local_max = -Float::INFINITY
  index_of_max = -1

  while right_ptr < nums.length
    if index_of_max >= left_ptr
      if nums[right_ptr] > local_max
        local_max = nums[right_ptr]
      end
    else
      window = nums[left_ptr..right_ptr]
      local_max = window.max
      index_of_max = window.find_index(local_max)
    end
    res << local_max
    right_ptr += 1
    left_ptr += 1
  end

  return res
end


#################### dynamic programming approach #####################

# The idea is to split an input array into blocks of k elements. The last block could contain less elements if n % k != 0.
# The pipes represent the boundaries of the blocks.
# We want to precompute 2 arrays, left and right such that
# left[i] represents the local maximums upto a block going from left to right
# And right[j] represents the local maximums upto a block going from right to left.

# For example :

# 1,3,-1 | -3,5,3 | 6,7
     # i      j

# Left = 1,3,3,-3,5,5,6,7
# Right = 3,3,-1,5,5,3,7,7

# Let's consider a sliding window from index i to index j.
# Element right[i] is a maximum element for window elements in the leftside block
# And element left[j] is a maximum element for window elements in the rightside block.
# Hence the maximum element in the sliding window is max(right[i], left[j])

# The algorithm is quite straightforward :
# 1. Iterate along the array in the direction left -> right and build an array left.
# 2. Iterate along the array in the direction left <- right and build an array right.
# 3. Build an output array as max(right[i], left[i + k - 1]) for i in range 0 to (nums.length - k)th element.


# @param {Integer[]} nums
# @param {Integer} k
# @return {Integer[]}
def max_sliding_window(nums, k)
  return [] if nums.length == 0
  return [] if k == 0
  return nums if nums.length == 1

  res = []
  left = Array.new(nums.length, 0)
  right = Array.new(nums.length, 0)
  left[0] = nums[0]
  right[nums.length - 1] = nums[nums.length - 1]

  (1..(nums.length-1)).each do |i|
    if i % k == 0
      left[i] = nums[i]
    else
      left[i] = [left[i-1], nums[i]].max
    end

    j = (nums.length - 1) - i
    if (j+1) % k == 0
      right[j] = nums[j]
    else
      right[j] = [right[j+1], nums[j]].max
    end
  end

  (0..(nums.length - k)).each do |m|
    res << [left[m+k-1], right[m]].max
  end

  return res
end
