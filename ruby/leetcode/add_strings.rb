# @param {String} num1
# @param {String} num2
# @return {String}
def add_strings(num1, num2)
  i = num1.length - 1
  j = num2.length - 1

  res = 0
  carry = 0
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
