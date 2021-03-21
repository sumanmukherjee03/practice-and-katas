# Given an array of meeting time intervals intervals where intervals[i] = [starti, endi], return the minimum number of conference rooms required.

# Example 1:

# Input: intervals = [[0,30],[5,10],[15,20]]
# Output: 2
# Example 2:

# Input: intervals = [[7,10],[2,4]]
# Output: 1

# Constraints:

# 1 <= intervals.length <= 104
# 0 <= starti < endi <= 106


# Algorithm

# Separate out the start times and the end times in their separate arrays.
# Sort the start times and the end times separately. Note that this will mess up the original correspondence of start times and end times. They will be treated individually now.
# We consider two pointers: s_ptr and e_ptr which refer to start pointer and end pointer. The start pointer simply iterates over all the meetings and the end pointer helps us track if a meeting has ended and if we can reuse a room.
# When considering a specific meeting pointed to by s_ptr, we check if this start timing is greater than the meeting pointed to by e_ptr. If this is the case then that would mean some meeting has ended by the time the meeting at s_ptr had to start. So we can reuse one of the rooms. Otherwise, we have to allocate a new room.
# If a meeting has indeed ended i.e. if start[s_ptr] >= end[e_ptr], then we increment e_ptr.
# Repeat this process until s_ptr processes all of the meetings.


# @param {Integer[][]} intervals
# @return {Integer}
def min_meeting_rooms(intervals)
  return 0 if intervals.empty?

  sorted_starts = intervals.map(&:first).sort # Get a sorted list of the start times
  sorted_ends = intervals.map(&:last).sort # Get a list of the sorted end times

  i = 0
  j = 0
  rooms = 0

  while i < sorted_starts.length do
    # If there is a meeting that has ended by the time the meeting at start pointer starts
    if sorted_starts[i] >= sorted_ends[j]
      rooms -= 1
      j += 1
    end

    # We do this irrespective of whether a room frees up or not.
    # If a room got free, then this rooms += 1 wouldn't have any effect. rooms would
    # remain the same in that case. If no room was free, then this would increase rooms
    rooms += 1
    i += 1
  end

  return rooms
end
