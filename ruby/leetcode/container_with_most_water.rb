# @param {Integer[]} height
# @return {Integer}
def max_area(height)
    max = 0

    get_area = lambda do |idx1, idx2|
        a = [height[idx1], height[idx2]].min * (idx2 - idx1)
        puts "#{idx1}, #{idx2}, #{a}"
        max = a if a > max
        get_area.call(idx1+1, idx2) if idx1 + 1 < idx2 && height[idx1] < height[idx2]
        get_area.call(idx1, idx2-1) if idx2 - 1 > idx1 && height[idx1] >= height[idx2]
    end

    get_area.call(0, height.length - 1)

    return max
end

