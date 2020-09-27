# @param {Integer[]} nums
# @param {Integer} target
# @return {Integer[]}
def two_sum(nums, target)
	i = 0
	res = []
	temp = {} # remainder => index of value
	while i < nums.length
		if temp[nums[i]]
			res = [i, temp[nums[i]]]
			break
		else
			temp[target - nums[i]] = i
		end
		i += 1
	end
	return res
end
