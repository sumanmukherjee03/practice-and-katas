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
