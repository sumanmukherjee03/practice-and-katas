# @param {Integer[][]} intervals
# @return {Integer}
def min_meeting_rooms(intervals)
  return 0 if intervals.empty?

  sorted_starts = intervals.map(&:first).sort
  sorted_ends = intervals.map(&:last).sort

  i = 0
  j = 0
  rooms = 0

  while i < sorted_starts.length do
    if sorted_starts[i] < sorted_ends[j]
      rooms += 1
      i += 1
    else
      i += 1
      j += 1
    end
  end

  return rooms
end
