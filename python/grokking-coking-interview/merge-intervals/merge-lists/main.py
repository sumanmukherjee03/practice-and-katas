from __future__ import print_function

class Interval:
    def __init__(self, start, end):
        self.start = start
        self.end = end

    def print_interval(self):
        print("[" + str(self.start) + ", " + str(self.end) + "]", end='')

def describe():
    desc = """
Problem : Given a list of intervals, merge all the overlapping intervals to produce a list that has only mutually exclusive intervals.
Example :
    Intervals: [[1,4], [2,5], [7,9]]
    Output: [[1,5], [7,9]]
    Explanation: Since the first two intervals [1,4] and [2,5] overlap, we merged them into one [1,5].
    """
    print(desc)

def merge(intervals):
    if len(intervals) < 2:
        return intervals
    #  Sort the interval based on the start time. This is O(nlog(n))
    intervals.sort(key=lambda x: x.start)
    merged = []
    current = intervals[0]
    for i in range(1, len(intervals)):
        interval = intervals[i]
        # Since the intervals are sorted by start, the only 2 cases we consider are
        #  <---a---->
        #         <----b---->
        #
        #  AND
        #
        #  <---a----> <---b---->
        if current.end >= interval.start:
            current.end = max(current.end, interval.end)
        else:
            merged.append(current)
            current = interval
    merged.append(current) # Dont forget to append the final interval
    return merged

def main():
    describe()

    print("Merged intervals: ", end='')
    for i in merge([Interval(1, 4), Interval(2, 5), Interval(7, 9)]):
        i.print_interval()
    print()

    print("Merged intervals: ", end='')
    for i in merge([Interval(6, 7), Interval(2, 4), Interval(5, 9)]):
        i.print_interval()
    print()

    print("Merged intervals: ", end='')
    for i in merge([Interval(1, 4), Interval(2, 6), Interval(3, 5)]):
        i.print_interval()
    print()

main()
