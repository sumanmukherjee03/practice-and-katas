# Given an array of intervals where intervals[i] = [starti, endi], merge all overlapping intervals,
# and return an array of the non-overlapping intervals that cover all the intervals in the input.
#
#  Example 1:
#
#  Input: intervals = [[1,3],[2,6],[8,10],[15,18]]
#  Output: [[1,6],[8,10],[15,18]]
#  Explanation: Since intervals [1,3] and [2,6] overlaps, merge them into [1,6].
#  Example 2:
#
#  Input: intervals = [[1,4],[4,5]]
#  Output: [[1,5]]
#  Explanation: Intervals [1,4] and [4,5] are considered overlapping.
#
#   Constraints:
#
#   1 <= intervals.length <= 104
#   intervals[i].length == 2
#   0 <= starti <= endi <= 104

# @param {Integer[][]} intervals
# @return {Integer[][]}
def merge(intervals)
  res = []

  # Sort the intervals based on the start of the interval. This O(nlogn)
  intervals.sort! {|x, y| x.first <=> y.first}

  # Go through the sorted list of intervals and look at the next interval and see
  # if the current interval overlaps with the previous interval.
  # As in, does the first element of the current set lie within the previous set.
  # If that overlap is there, then the sets should be merged
  intervals.each do |set|
    # If the resulting list is empty then insert the first item into the resulting set
    if res.length == 0
      res << set
      next
    end

    last_set = res.last
    begining_of_last_merged_set = last_set.first
    end_of_last_merged_set = last_set.last
    begining_of_current_set = set.first
    end_of_current_set = set.last

    if begining_of_current_set <= end_of_last_merged_set
      start = begining_of_last_merged_set
      if end_of_last_merged_set <= end_of_current_set
        stop = end_of_current_set
      else
        stop = end_of_last_merged_set
      end
      res[-1] = [start, stop]
    else
      res << set
    end
  end

  res
end
