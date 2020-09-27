# @param {Integer} x
# @return {Boolean}
def is_palindrome(x)
    return false if x < 0
    y = reverse(x)
    return x == y
end

def reverse(x)
    arr = []

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

    return num
end
