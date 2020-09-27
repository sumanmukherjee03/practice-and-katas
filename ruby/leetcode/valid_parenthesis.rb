# @param {String} s
# @return {Boolean}
def is_valid(s)
    chars = s.chars
    valid = true
    stack = []

    chars.each do |c|
        case c
        when ')'
            if stack.last == '('
                stack.pop
            else
                valid = false
                break
            end
        when '}'
            if stack.last == '{'
                stack.pop
            else
                valid = false
                break
            end
        when ']'
            if stack.last == '['
                stack.pop
            else
                valid = false
                break
            end
        else
            stack << c
        end
    end

    return valid && stack.empty?
end
