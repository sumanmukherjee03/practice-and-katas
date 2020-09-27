# @param {String} s
# @return {String}
def longest_palindrome(s)
  arr = s.chars
  dp = init_dp_array_with_diag_set(arr.length)
  res = ""

  arr.each_with_index do |c, j|
    i = 0
    while i <= j do
      if j - i > 1
        if dp[i+1][j-1] == 1 && arr[i] == arr[j]
          dp[i][j] = 1
        end
      else
        if arr[i] == arr[j]
          dp[i][j] = 1
        end
      end

      if dp[i][j] == 1 && res.length < j - i + 1
        res = arr[i..j].join("")
      end

      i += 1
    end
  end
  return res
end

def init_dp_array(dim)
  return Array.new(dim){Array.new(dim, 0)} # initialize empty 2d array for dynamic programming
end

def init_dp_array_with_diag_set(dim)
  dp = init_dp_array(dim)
  x = 0
  while x < dim do
    dp[x][x] = 1
    x += 1
  end
  return dp
end
