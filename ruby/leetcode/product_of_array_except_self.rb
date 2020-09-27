# @param {Integer[]} nums
# @return {Integer[]}
def product_except_self(nums)
    product_before = Array.new(nums.length, 1)
    product_after = Array.new(nums.length, 1)

    i = 1
    while i < nums.length do
        product_before[i] = product_before[i-1] * nums[i-1]
        i += 1
    end

    product_after[nums.length - 1] = 1
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
