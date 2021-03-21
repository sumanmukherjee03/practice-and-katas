# Given a string s, return the longest palindromic substring in s.

# Example 1:

# Input: s = "babad"
# Output: "bab"
# Note: "aba" is also a valid answer.
# Example 2:

# Input: s = "cbbd"
# Output: "bb"
# Example 3:

# Input: s = "a"
# Output: "a"
# Example 4:

# Input: s = "ac"
# Output: "a"

# Constraints:

# 1 <= s.length <= 1000
# s consist of only digits and English letters (lower-case and/or upper-case)


# Algorithm
# we first observe how we can avoid unnecessary re-computation while validating palindromes.
# Consider the case "ababa". If we already knew that "bab" is a palindrome,
# it is obvious that "ababa" must be a palindrome since the two left and right end letters are the same

# We define P(i,j) as following :
  # P(i, j) = true if the substring Si..Sj is a palindrome  else false
# P(i,j)=(P(i+1,j−1) and Si == Sj

# Base cases :
# P(i,i)=true
# P(i,i+1)= (Si == Sj)
# This yields a straight forward DP solution, which we first initialize the one and two letters palindromes,
# and work our way up finding all three letters palindromes, and so on

# Time complexity : O(n^2)
# Space complexity : O(n ^ 2)


# @param {String} s
# @return {String}
def longest_palindrome(s)
  arr = s.chars

  # All chars by themselves are palindromes
  dp = init_dp_array_with_diag_set(arr.length)

  res = ""

  arr.each_with_index do |c, j|
    i = 0
    while i <= j do
      # Consider case for substrings bigger than 2 chars
      if j - i > 1
        # If substring from position i+1 to j-1 is a palindrome and char at i and char at j are same
        # then substring from position i to position j is also a palindrome
        if dp[i+1][j-1] == 1 && arr[i] == arr[j]
          dp[i][j] = 1
        end
      else
        # This is the case for substrings with 2 chars or 1 char
        # If both the chars in that substring are the same then it is a palindrome
        if arr[i] == arr[j]
          dp[i][j] = 1
        end
      end

      if dp[i][j] == 1 && res.length < (j - i) + 1
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
