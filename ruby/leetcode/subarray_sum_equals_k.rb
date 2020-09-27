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
