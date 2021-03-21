# Given a string s, find the length of the longest substring without repeating characters.

# Example 1:

# Input: s = "abcabcbb"
# Output: 3
# Explanation: The answer is "abc", with the length of 3.
# Example 2:

# Input: s = "bbbbb"
# Output: 1
# Explanation: The answer is "b", with the length of 1.
# Example 3:

# Input: s = "pwwkew"
# Output: 3
# Explanation: The answer is "wke", with the length of 3.
# Notice that the answer must be a substring, "pwke" is a subsequence and not a substring.
# Example 4:

# Input: s = ""
# Output: 0

# Constraints:

# 0 <= s.length <= 5 * 104
# s consists of English letters, digits, symbols and spaces.


# @param {String} s
# @return {Integer}
def length_of_longest_substring(s)
  arr = s.chars
  h = {}
  i = 0

  # From a particular position in the string chars scan the rest of the chars
  # to see whats the max length of substring without repeating chars to be found
  max_subsrt_len_found = 0
  while i < arr.length do
    arr[i..-1].each do |c|
      if h.has_key?(c)
        break
      else
        h[c] = 1
      end
    end
    max_subsrt_len_found = h.keys.length if h.keys.length > max_subsrt_len_found
    h = {}
    i += 1
  end
  return res
end
