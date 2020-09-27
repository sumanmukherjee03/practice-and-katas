# @param {Integer} n
# @return {String[]}
def generate_parenthesis(n)
    chars = ['(', ')']

    res = []
    visited = {}

    solve = lambda do |ith_pos,possible_solution|
      return if possible_solution.length > (chars.length * n)
      return if ith_pos >= (chars.length * n) || ith_pos < 0

      # For the ith_pos try all the chars
      j = 0
      while j < chars.length do
        c = chars[j]
        solution = possible_solution.clone
        solution[ith_pos] = c
        if solution.length == (chars.length * n)
          str = solution.join('')
          break if visited.has_key?(str)

          if is_paranthesis_valid?(str)
            visited[str] = 1
            res << str
            break
          else
            visited[str] = 0
            j += 1
          end
        else
          solve.call(ith_pos+1,solution)
          j += 1
        end
      end
    end
    solve.call(0,[])

    return res
end

def is_paranthesis_valid?(s)
  chars = s.chars
  stack = []
  mismatched_closing_paranthesis_found = false

  i = 0
  while i < chars.length do
    if chars[i] == '('
      stack.push(chars[i])
    else
      if !stack.empty?
        stack.pop
      else
        mismatched_closing_paranthesis_found = true
      end
    end
    i += 1
  end

  return stack.empty? && !mismatched_closing_paranthesis_found
end
