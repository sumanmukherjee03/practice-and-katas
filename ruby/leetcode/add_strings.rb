# Given two non-negative integers num1 and num2 represented as string, return the sum of num1 and num2.

# Note:

# The length of both num1 and num2 is < 5100.
# Both num1 and num2 contains only digits 0-9.
# Both num1 and num2 does not contain any leading zero.
# You must not use any built-in BigInteger library or convert the inputs to integer directly


# @param {String} num1
# @param {String} num2
# @return {String}
def add_strings(num1, num2)
  i = num1.length - 1
  j = num2.length - 1

  res = 0
  while i >= 0 || j >= 0 do
    if i >= 0 && j >= 0
      res = res + ((num1[i].to_i + num2[j].to_i) * (10 ** (num1.length - i - 1)))
      i -= 1
      j -= 1
    elsif i >= 0 && j < 0
      res = res + ((num1[i].to_i) * (10 ** (num1.length - i - 1)))
      i -= 1
    elsif i < 0 && j >= 0
      res = res + ((num2[j].to_i) * (10 ** (num2.length - j - 1)))
      j -= 1
    end
  end

  return res.to_s
end
