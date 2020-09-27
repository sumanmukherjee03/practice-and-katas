# @param {String} s
# @return {Integer}
def length_of_longest_substring(s)
  arr = s.chars
  h = {}
  i = 0
  res = 0
  while i < arr.length do
    arr[i..-1].each do |c|
      if h.has_key?(c)
        res = h.keys.length if h.keys.length > res
        break
      else
        h[c] = 1
      end
    end
    res = h.keys.length if h.keys.length > res
    h = {}
    i += 1
  end
  return res
end
