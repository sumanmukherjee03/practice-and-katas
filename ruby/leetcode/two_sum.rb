# Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.
#
# You may assume that each input would have exactly one solution, and you may not use the same element twice.
#
# You can return the answer in any order.
#
#  Example 1:
#
#  Input: nums = [2,7,11,15], target = 9
#  Output: [0,1]
#  Output: Because nums[0] + nums[1] == 9, we return [0, 1].
#  Example 2:
#
#  Input: nums = [3,2,4], target = 6
#  Output: [1,2]
#  Example 3:
#
#  Input: nums = [3,3], target = 6
#  Output: [0,1]
#
#
#   Constraints:
#
#   2 <= nums.length <= 103
#   -109 <= nums[i] <= 109
#   -109 <= target <= 109
#   Only one valid answer exists.
#
#
# @param {Integer[]} nums
# @param {Integer} target
# @return {Integer[]}

def two_sum(nums, target)
	i = 0
	res = []
	temp = {} # k:v would be : remainder => index of value
	while i < nums.length
    # If there is a key in temp that matches the number at index i
		if temp[nums[i]]
			res = [i, temp[nums[i]]]
			break
    # Keep filling the remainders as key and value as the index of the element that is currently considered
		else
			temp[target - nums[i]] = i
		end
		i += 1
	end
	return res
end
