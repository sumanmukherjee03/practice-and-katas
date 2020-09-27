# @param {Integer[]} nums
# @return {Integer[][]}
def three_sum(nums)
  solution = []

    h = {}
    x = []

    i = 0
    while i < nums.length do
      h[0 - nums[i]] = [i]
      i += 1
    end
    puts "h : #{h}"

    h.each do |k, v|
      puts "Key : #{k}, Value: #{v}"

      h1 = {}
      j = 0
      found = false
      index = v.first
      while j < nums.length do
        if j != index
          h1[k - nums[j]] = j
        end
        j += 1
      end

      puts "h1 : #{h1}"

      m = 0
      while m < nums.length do
        if m != index
          if h1.has_key?(nums[m])
            found = true
            break
          end
        end
      end

      if found
        s = [nums[index], nums[h1[nums[m]]], nums[m]].sort
        solution << s
        solution.uniq!
        puts "solution : #{solution}"
      end
    end

    return solution
end
