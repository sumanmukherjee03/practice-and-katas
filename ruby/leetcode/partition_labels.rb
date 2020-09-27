# @param {String} s
# @return {Integer[]}
def partition_labels(s)
  # "ababcbacadefegdehijhklij"
  chars = s.chars
  sets = []
  h = {} # key is char, value is index of the set

  i = 0
  while i < chars.length do
    if i-1 >= 0
      previous_c = chars[i-1]
      previous_c_idx = h[previous_c]
    end

    c = chars[i]

    if h.has_key?(c)
      c_idx = h[c]
      if previous_c && previous_c_idx > c_idx
        sets[c_idx] = sets[c_idx..previous_c_idx].flatten
        sets[c_idx] << c
        sets = sets[0..c_idx]
        (c_idx+1..previous_c_idx).each do |idx|
          h.select{|x, y| y == idx }.map {|z| z.first}.each do |key|
            h[key] = c_idx
          end
        end
      else
        sets[c_idx] << c
      end
    else
      sets << [c]
      h[c] = sets.length - 1
    end

    i += 1
  end

  return sets.map {|s| s.length}
end
