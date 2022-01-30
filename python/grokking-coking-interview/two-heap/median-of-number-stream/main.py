from heapq import *

def describe():
    desc = """
Problem : Design a class to calculate the median of a number stream. The class should have the following two methods:

    insertNum(int num): stores the number in the class
    findMedian(): returns the median of all numbers inserted in the class

If the count of numbers inserted in the class is even, the median will be the average of the middle two numbers.
Remember that the median is the middle value in an ordered integer list.

-------------
    """
    print(desc)

#  As we know, the median is the middle value in an ordered integer list.
#  So a brute force solution could be to maintain a sorted list of all numbers inserted in the class
#  so that we can efficiently return the median whenever required.
#  Inserting a number in a sorted list will take O(N) time if there are N numbers in the list.
#  This insertion will be similar to the Insertion sort. However there is a better approach.
#  Assume x is the median of a list. This means that half of the numbers in the list will be smaller than (or equal to) x
#  and half will be greater than (or equal to) x. This leads us to an approach where we can divide the list into two halves.
#  One half to store all the smaller numbers (let's call it smallNumList) and one half to store the larger numbers (let's call it largNumList).
#  The median of all the numbers will either be the largest number in the smallNumList or the smallest number in the largNumList.
#  If the total number of elements is even, the median will be the average of these two numbers.
#  The best data structure to do this is a Heap.

#  The process could go like this.
#  Maintain 2 heaps - left part is a max heap because we want to find the largest number from the left.
#  And the right part can be a minheap because we want to find the smallest number from the right.
#  After insertion we try to balance out the number of elements in the 2 heaps so that they are equal.
#  If there are odd number of elements in the heap then we decide to have more numbers in the max-heap , ie the left part.
#  If the min-heap ends up having more elements than the max-heap after insertion, we can move elements from minheap to max heap
#  because we have decided to keep more elements in the maxheap than the minheap.

class MedianOfAStream():
    def __init__(self):
        # Start with simple arrays as the collecting data structures
        #  And as you insert elements, it will heapify and form the heaps
        self.maxheap = []
        self.minheap = []

    #  heapq by default maintains a minheap. So, turn that into a maxheap just store negative of the numbers. Thats the trick
    #  When withdrawing from the maxheap always negate the number again, so that you get back the right value
    #  And the top of the heap is always the 0th element
    #  Time complexity of insertion in heap is O(logn)
    def insert_num(self, num):
        if not self.maxheap or -self.maxheap[0] >= num:
            heappush(self.maxheap, -num)
        else:
            heappush(self.minheap, num)

        if len(self.maxheap) > len(self.minheap) + 1:
            heappush(self.minheap, -heappop(self.maxheap))
        elif len(self.maxheap) < len(self.minheap):
            heappush(self.maxheap, -heappop(self.minheap))

    def find_median(self):
        if len(self.maxheap) == len(self.minheap):
            return (-self.maxheap[0] + self.minheap[0])/2.0
        return -self.maxheap[0]


def main():
    describe()
    medianOfAStream = MedianOfAStream()
    medianOfAStream.insert_num(3)
    medianOfAStream.insert_num(1)
    print("The median is: " + str(medianOfAStream.find_median()))
    medianOfAStream.insert_num(5)
    print("The median is: " + str(medianOfAStream.find_median()))
    medianOfAStream.insert_num(4)
    print("The median is: " + str(medianOfAStream.find_median()))

main()
