# Trapping Rain Water
# Given n non-negative integers representing an elevation map where the width of each bar is 1, compute how much water it can trap after raining.

# Example 1:


# Input: height = [0,1,0,2,1,0,1,3,2,1,2,1]
# Output: 6
# Explanation: The above elevation map (black section) is represented by array [0,1,0,2,1,0,1,3,2,1,2,1]. In this case, 6 units of rain water (blue section) are being trapped.
# Example 2:

# Input: height = [4,2,0,3,2,5]
# Output: 9


# Constraints:

# n == height.length
# 0 <= n <= 3 * 104
# 0 <= height[i] <= 105


# @param {Integer[]} height
# @return {Integer}
def trap(height)
    return 0 if height.length == 0

    left_max = Array.new(height.length, 0)
    right_max = Array.new(height.length, 0)

    # Do one iteration and pre compute the max lefts for a given bar
    i = 0
    l_max = 0
    while i < left_max.length do
        l_max = i if height[i] >= height[l_max]
        left_max[i] = l_max
        i += 1
    end

    # Do one iteration and pre compute the max rights for a given bar
    j = right_max.length - 1
    r_max = right_max.length - 1
    while j >= 0 do
        r_max = j if height[j] >= height[r_max]
        right_max[j] = r_max
        j -= 1
    end

    # Calculate area per position
    k = 0
    area = 0
    while k < height.length do
        a = ([height[left_max[k]],height[right_max[k]]].min - height[k]) * 1
        area += a
        k += 1
    end

    return area
end
