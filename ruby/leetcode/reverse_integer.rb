# Given a signed 32-bit integer x, return x with its digits reversed. If reversing x causes the value to go outside the signed 32-bit integer range [-231, 231 - 1], then return 0.

# Assume the environment does not allow you to store 64-bit integers (signed or unsigned).

# Example 1:

# Input: x = 123
# Output: 321
# Example 2:

# Input: x = -123
# Output: -321
# Example 3:

# Input: x = 120
# Output: 21
# Example 4:

# Input: x = 0
# Output: 0

# Constraints:

# -231 <= x <= 231 - 1

# @param {Integer} x
# @return {Integer}
def reverse(x)
    arr = []
    is_negative = x < 0 ? true : false

    x = x.abs

    while x > 0 do
        arr << x % 10
        x = x/10
    end

    num = 0
    i = 0
    while i < arr.length do
        d = (arr[i] * (10 ** (arr.length - i - 1)))
        num += d
        i += 1
    end

    return 0 if num > 2147483647 || num < -2147483648
    return is_negative ? 0 - num : num
end
