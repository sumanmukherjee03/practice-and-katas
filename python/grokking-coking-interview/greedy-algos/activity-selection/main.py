from __future__ import print_function

def describe():
    desc = """
Problem : You are given n activities with their start and finish times.
          Select the maximum number of activities that can be performed by a single person, assuming that a person can only work on a single activity at a time.
Example :
    Input : start = [10, 12, 20], finish = [20, 25, 30]
    Output : Max set of activities that can be executed is [0, 2]. These are the indices of the activities in start/finsh.

    Input : start = [1, 3, 0, 5, 8, 5], finish = [2, 4, 6, 7, 9, 9]
    Output : Max set of activities that can be executed is [0,1,3,4]. These are the indices of the activities in start/finsh.


--------------------
    """
    print(desc)

#  Algo explanation :
#      The greedy choice is to pick the next activity whose finish time is the closest among the remaining activities.
#      And the start time is after the finish time of the previously selected activity.
#
#  So, sort the activities according to their finish time.
#  Select the first activity from the sorted array and print it.
#  Repeat this for the rest of the sorted array of activities by finish time
#     - If the start time of the next activity in the sorted list is greater than the finish time of the previously selected activity then take that and print it

#  The first solution we are exploring are of activities that are already sorted by their finish time.
#  Assuming that the activities are already sorted by their finish time.
def printMaxActivities(start, finish):
    n = len(finish)
    print("The following activites were selected : ", end='')
    i = 0
    print(i, end=', ')

    for j in range(1, n):
        if start[j] >= finish[i]:
            print(j, end=', ')
            i = j





#  The second example is of activities that are not sorted by their finish time.
def maxActivities(activities):
    selected = []
    activities.sort(key = lambda x : x[1]) # Sort the activities by their finish time

    i = 0
    selected.append(activities[i])

    for j in range(1, len(activities)):
        if activities[j][0] >= activities[i][1]:
            selected.append(activities[j])
            i = j

    return selected



def main():
    describe()

    start = [10, 12, 20]
    finish = [20, 25, 30]
    print("Input :")
    print(str(start))
    print(str(finish))
    printMaxActivities(start, finish)
    print("\n\n")

    activities = [[0,6], [3,4], [5,7], [1,2], [5,9], [8,9]]
    print("Input :")
    print(str(activities))
    result = maxActivities(activities)
    print(str(result))
    print("\n\n")

main()
