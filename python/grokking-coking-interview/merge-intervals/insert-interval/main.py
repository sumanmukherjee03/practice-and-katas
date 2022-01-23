def describe():
    desc = """
Problem : Given a list of non-overlapping intervals sorted by their start time, insert a given interval at the correct position and merge all necessary intervals to produce a list that has only mutually exclusive intervals.
Example :
    Input: Intervals=[[1,3], [5,7], [8,12]], New Interval=[4,6]
    Output: [[1,3], [4,7], [8,12]]
    Explanation: After insertion, since [4,6] overlaps with [5,7], we merged them into one [4,7].
--------------
    """
    print(desc)


#  Since this list is already sorted this will get done in O(n) time
def insert(intervals, new_interval):
  merged = []
  i, start, end = 0, 0, 1

  #  skip (and add to output) all intervals that come before the 'new_interval'
  #  ie whose end is before the start of the new interval
  while i < len(intervals) and intervals[i][end] < new_interval[start]:
    merged.append(intervals[i])
    i += 1

  #  merge all intervals that overlap with 'new_interval'
  #  overlap will be those intervals whose start is between the start and stop of the new interval
  while i < len(intervals) and intervals[i][start] <= new_interval[end]:
    new_interval[start] = min(intervals[i][start], new_interval[start])
    new_interval[end] = max(intervals[i][end], new_interval[end])
    i += 1

  # insert the new_interval
  merged.append(new_interval)

  # add all the remaining intervals to the output
  while i < len(intervals):
    merged.append(intervals[i])
    i += 1
  return merged

def main():
    describe()
    print("Intervals after inserting the new interval: " + str(insert([[1, 3], [5, 7], [8, 12]], [4, 6])))
    print("Intervals after inserting the new interval: " + str(insert([[1, 3], [5, 7], [8, 12]], [4, 10])))
    print("Intervals after inserting the new interval: " + str(insert([[2, 3], [5, 7]], [1, 4])))

main()
