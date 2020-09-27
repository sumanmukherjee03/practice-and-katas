# @param {String} s
# @return {String}
def min_remove_to_make_valid(s)
    chars = s.chars
    stack = []
    marked_for_deletion = []

    i = 0
    while i < chars.length do
      if chars[i] == '('
        stack.push([i, chars[i]])
      elsif chars[i] == ')'
        if stack.empty?
          marked_for_deletion << i
        else
          stack.pop
        end
      end
      i += 1
    end

    marked_for_deletion.concat(stack.map {|x| x[0]}.flatten) unless stack.empty?

    unless marked_for_deletion.empty?
      j = 0
      while j < marked_for_deletion.length do
        idx = marked_for_deletion[j]
        chars[idx] = '$'
        j += 1
      end
    end

    return chars.join('').gsub('$', '')
end
