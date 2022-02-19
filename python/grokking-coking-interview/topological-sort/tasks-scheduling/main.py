from collections import deque

def describe():
    desc = """
Problem : There are N tasks labelled 0 to N-1. Each task has some pre-requisite task that needs to be completed
          before it can be started. Given the number of tasks and a list of prerequisite pairs, find out if it is
          possible to schedule all the tasks.
Example :
        Input: Tasks=3, Prerequisites=[0, 1], [1, 2]
        Output: true
        Explanation: To execute task 1, task 0 needs to finish first. Similarly, task 1 needs
        to finish before 2 can be scheduled. One possible scheduling of tasks is: [0, 1, 2]

        Input: Tasks=3, Prerequisites=[0, 1], [1, 2], [2, 0]
        Output: false
        Explanation: The tasks have a cyclic dependency, therefore they cannot be scheduled.

        Input: Tasks=6, Prerequisites=[2, 5], [0, 5], [0, 4], [1, 4], [3, 2], [1, 3]
        Output: true
        Explanation: A possible scheduling of tasks is: [0 1 4 3 2 5]

--------------------
    """
    print(desc)

#  Time complexity of this is O(V+E)
def is_scheduling_possible(numOfTasks, prerequisites):
    jobSchedule = []

    inDegrees = {i: 0 for i in range(0, numOfTasks)}
    adjacencyList = {i: [] for i in range(0, numOfTasks)}

    for prereq in prerequisites:
        src, dest = prereq
        adjacencyList[src].append(dest)
        inDegrees[dest] += 1

    q = deque()
    for t in inDegrees:
        if inDegrees[t] == 0:
            q.append(t)

    while q:
        n = len(q)
        for _ in range(n):
            task = q.popleft()
            jobSchedule.append(task)
            for t in adjacencyList[task]:
                inDegrees[t] -= 1
                if inDegrees[t] == 0:
                    q.append(t)

    print("Possible job schedule is : " + str(jobSchedule))

    #  IMPORTANT NOTE : If there is a cycle then the length of the jobSchedule array will be less than numOfTasks
    #  This condition can be used to determine if a directed graph has a cycle or not.
    #  This is because where there is a cycle, we wouldnt be able to find any sources, so the queue would be empty.
    #  We would only be able to sort the jobs before the cycle starts.
    if len(jobSchedule) != numOfTasks:
        return False

    return True

def main():
    describe()

    tasks = 3
    prerequisites = [[0, 1], [1, 2]]
    print("Input : " + str(tasks) + ", " + str(prerequisites))
    print("Topological sort: " + str(is_scheduling_possible(tasks, prerequisites)))
    print("\n\n")

    tasks = 3
    prerequisites = [[0, 1], [1, 2], [2, 0]]
    print("Input : " + str(tasks) + ", " + str(prerequisites))
    print("Topological sort: " + str(is_scheduling_possible(tasks, prerequisites)))
    print("\n\n")

    tasks = 6
    prerequisites = [[0, 4], [1, 4], [3, 2], [1, 3]]
    print("Input : " + str(tasks) + ", " + str(prerequisites))
    print("Topological sort: " + str(is_scheduling_possible(tasks, prerequisites)))
    print("\n\n")

main()
