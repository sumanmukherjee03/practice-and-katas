# @param {String[]} strs
# @return {String[][]}
def group_anagrams(strs)
  hashes = {}

  i = 0
  while i < strs.length do
    chars = strs[i].chars
    h = {}
    j = 0
    while j < chars.length do
      if h.has_key?(chars[j])
        h[chars[j]] += 1
      else
        h[chars[j]] = 1
      end
      j += 1
    end

    s = h.sort_by(&:first).inject("") {|str, arr| str += arr[0] + arr[1].to_s; str}

    if hashes.has_key?(s)
      hashes[s] << strs[i]
    else
      hashes[s] = [strs[i]]
    end

    i += 1
  end

  return hashes.values
end
