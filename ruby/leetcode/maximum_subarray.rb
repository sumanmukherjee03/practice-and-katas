# @param {Integer[]} nums
# @return {Integer}
def max_sub_array(nums)
    i = 0
    max_till_current = 0
    max = nil
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

    j = 0
    while j < nums.length do
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
