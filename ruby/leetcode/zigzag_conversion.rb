# @param {String} s
# @param {Integer} num_rows
# @return {String}
def convert(s, num_rows)
    if num_rows == 1 || num_rows >= s.length
        return s
    end

    chars = s.chars

    no_diags = chars.length/num_rows - 1
    no_cols = chars.length/num_rows + no_diags
    myarr = Array.new(num_rows){Array.new(no_cols, " ")}

    i = 0
    col = 0
    row = 0
    is_top_to_bottom = true
    while i < chars.length do
        if !is_top_to_bottom
            if row+1 < num_rows && myarr[row+1][col-1] != " "
                myarr[row][col] = chars[i]
                i += 1
                col += 1
            end
            row -= 1
        else
            myarr[row][col] = chars[i]
            i += 1
            row += 1
        end

        if row == 0
            is_top_to_bottom = true
        end

        if row == num_rows
            is_top_to_bottom = false
            row = num_rows - 1
            col += 1
        end
    end

    str = ""
    myarr.each do |row|
        str += row.join('').delete(' ')
    end

    return str
end
