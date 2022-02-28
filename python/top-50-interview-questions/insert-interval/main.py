def describe():
    desc = """
Problem : Given a new interval and an array of intervals, insert the new interval and merge if necessary.
          Note that the intervals in the array are non-overlapping and are sorted by their starting point.

Example :
    Input : intervals : [[1,3], [4,7], [8,10], [12,15], [16,17], [18,20], [21,25], [28,29]], new interval [9,18]
    Output : [[1,3], [4,7], [8,20], [21,25], [28,29]]
    """
    print(desc)




#  Time complexity is O(n)
def insertInterval(intervals, newInterval):
    if len(intervals) == 0:
        return newInterval

    # The intervals are sorted by their start index.
    #  So, first find the position where the new interval can be inserted and then insert it there.
    i = 0
    while newInterval[0] >= intervals[i][0] and i < len(intervals):
        i += 1
    intervals = intervals[0:i] + [newInterval] + intervals[i:]

    # Next if there are overlapping intervals created due to the insert, merge them
    k = 0
    while k < len(intervals):
        if k+1 < len(intervals) and intervals[k][1] >= intervals[k+1][0]:
            start = intervals[k][0]
            stop = max(intervals[k][1], intervals[k+1][1])
            intervals[k+1] = [start, stop]
            intervals[k] = []
        k += 1
    result = [interval for interval in intervals if len(interval) > 0]
    return result




#  Time complexity is O(n)
def insertIntervalOptimized(intervals, newInterval):
    output = []
    i = 0
    #  First find out where the new interval needs to be inserted.
    #  It will be inserted after the interval where the end of the previous interval is before the start of the new interval
    while i < len(intervals) and intervals[i][1] < newInterval[0]:
        output.append(intervals[i])
        i += 1
    #  Then merge the intervals from that point onwards with the new interval
    while i < len(intervals) and intervals[i][0] <= newInterval[1]:
        newInterval[0] = min(newInterval[0], intervals[i][0])
        newInterval[1] = max(newInterval[1], intervals[i][1])
        i += 1
    output.append(newInterval)
    while i < len(intervals):
        output.append(intervals[i])
        i += 1
    return output





def main():
    describe()

    intervals = [[1,3], [4,7], [8,10], [12,15], [16,17], [18,20], [21,25], [28,29]]
    newInterval = [9,18]
    print("Input : " + str(intervals) + ", " + str(newInterval))
    print("Output : " + str(insertInterval(intervals, newInterval)))

main()
