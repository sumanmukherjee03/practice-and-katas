from __future__ import print_function
import heapq

class Point(object):

    """Docstring for Point. """

    def __init__(self, x, y):
        self.x = x
        self.y = y

    #  This is necessary for comparison in maxheap
    #  Remember for maxheap we use a negative value. So, we are using the greater than operator in a lt function
    def __lt__(self, other):
        return self.dist_from_origin() > other.dist_from_origin()

    #  Ignoring the sqrt to calculate the distance here
    def dist_from_origin(self):
        return (self.x*self.x + self.y*self.y)

    def print_point(self):
        print("[" + str(self.x) + "," + str(self.y) + "]", end='')

def describe():
    desc = """
Problem : Given an array of 2D points in a plane, find the k closest points to origin.
Example:
    Input: points = [[1,2],[1,3]], K = 1
    Output: [[1,2]]
    Explanation: The Euclidean distance between (1, 2) and the origin is sqrt(5).
    The Euclidean distance between (1, 3) and the origin is sqrt(10).
    Since sqrt(5) < sqrt(10), therefore (1, 2) is closer to the origin.

    Input: point = [[1, 3], [3, 4], [2, -1]], K = 2
    Output: [[1, 3], [2, -1]]
    """
    print(desc)

def find_closest_points(points, k):
    maxheap = []
    for i in range(0, len(points)):
        if len(maxheap) < k:
            heapq.heappush(maxheap, points[i])
        else:
            if maxheap[0].dist_from_origin() > points[i].dist_from_origin():
                heapq.heappop(maxheap)
                heapq.heappush(maxheap, points[i])
    return maxheap


def main():
    describe()
    result = find_closest_points([Point(1, 3), Point(3, 4), Point(2, -1)], 2)
    print("Here are the k points closest the origin: ", end='')
    for point in result:
        point.print_point()

main()
