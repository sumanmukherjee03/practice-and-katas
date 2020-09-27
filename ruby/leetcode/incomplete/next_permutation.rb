# @param {Integer[]} nums
# @return {Void} Do not return anything, modify nums in-place instead.
def next_permutation(nums)
  return nums if nums.length <= 1

  result = nums.sort
  visited = []

  # Generate the current value
  orig_val = 0
  nums.reverse.each_with_index {|n,i| orig_val += n * (10**i)}

  min_next = Float::INFINITY

  solve = lambda do |ith_position, orig_solution|
    return if ith_position >= nums.length || ith_position < 0
    digits_remaining = nums.dup
    orig_solution.each do |n|
      idx = digits_remaining.find_index(n)
      digits_remaining.delete_at(idx)
    end

    # go over the remaining digits and place them in ith_position one by one
    j = 0
    while j < digits_remaining.length do
      possible_solution = orig_solution.dup
      possible_solution[ith_position] = digits_remaining[j]

      # Generate value for new solution
      possible_solution_val = 0
      possible_solution.reverse.each_with_index {|n,i| possible_solution_val += n * (10**i)}

      # if solution is of the same length as max, then you have reached a possible solution
      if possible_solution.length == nums.length
        visited << possible_solution_val
        # if new solution is greater than original value and less than the min next found earlier
        if possible_solution_val > orig_val && possible_solution_val < min_next
          min_next = possible_solution_val # make min next the value of the new solution
          result = possible_solution # make result the new possible solution
        else
          j += 1
        end
      else
        solve.call(ith_position+1,possible_solution)
        j += 1
      end
    end
  end

  solve.call(0, [])

  result.each_with_index {|n,i| nums[i] = n}
end
