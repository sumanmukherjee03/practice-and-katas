import heapq

def describe():
    desc = """
Problem : Given an array of numbers and a number k, find the median of all k sized subarrays or windows of that array.
Remember that the median is the middle value in an ordered integer list.

Example :
    Input: nums=[1, 2, -1, 3, 5], k = 2
    Output: [1.5, 0.5, 1.0, 4.0]
    Explaination :
        [1, 2] -> median is 1.5
        [2, -1] -> median is 1.5
        ... and so on

-----------------
    """
    print(desc)

class SlidingWindowMedian():
    def __init__(self):
        # Start with simple arrays as the collecting data structures
        #  And as you insert elements, it will heapify and form the heaps
        self.maxheap = [] # This represents the left partition because we want the max from this partition for median. It will also contain 1 more element than the right partition
        self.minheap = [] # This represents the right partition because we want the min from this partition for the median.

    def remove(self, heap, element):
        idx = heap.index(element)
        #  Delete the element by exchanging it with the last element of the array
        #  And then deleting from the last element of the array
        #  However, this means we will be losing the heap property because the previously placed last element is now at a wrong position
        #  So we will need to re-heapify the data structures
        heap[idx] = heap[-1]
        del heap[-1]
        if idx < len(heap):
            # We can use heapify to readjust the elements but that would be O(N),
            # instead, we will adjust only one element which will O(logN)
            heapq._siftup(heap, idx)
            heapq._siftdown(heap, 0, idx)

    def rebalance_heaps(self):
        # either both the heaps will have equal number of elements or max-heap will have
        # one more element than the min-heap
        if len(self.maxheap) > len(self.minheap) + 1:
            heapq.heappush(self.minheap, -heapq.heappop(self.maxheap))
        elif len(self.maxheap) < len(self.minheap):
            heapq.heappush(self.maxheap, -heapq.heappop(self.minheap))

    def find_sliding_window_median(self, nums, k):
        result = [0.0 for x in range(len(nums) - k + 1)]
        for i in range(0, len(nums)):
            if not self.maxheap or nums[i] <= -self.maxheap[0]:
                heapq.heappush(self.maxheap, -nums[i])
            else:
                heapq.heappush(self.minheap, nums[i])

            self.rebalance_heaps()

            #  If we have at least k elements in the sliding window
            #  Up until that hold off from calculating the median and just keep inserting elements into the max and min heaps
            #  When the iteration is at k elements and more, that's when the removal and insertion from the sliding window starts happening
            if i - k + 1 >= 0:
                # Calculate and add the median to the the result array
                if len(self.maxheap) == len(self.minheap):
                    # we have even number of elements, take the average of middle two elements
                    result[i - k + 1] = (-self.maxheap[0] + self.minheap[0]) / 2.0
                else:  # because max-heap will have one more element than the min-heap
                    result[i - k + 1] = -self.maxheap[0] / 1.0

                # remove the element going out of the sliding window
                elementToBeRemoved = nums[i - k + 1]

                #  If the element is less than the max in the min heap, that means it is in the left partion
                #  Otherwise in the right partition
                #  The insertion of new element in the sliding window has already happened at the top
                if elementToBeRemoved <= -self.maxheap[0]:
                    self.remove(self.maxheap, -elementToBeRemoved)
                else:
                    self.remove(self.minheap, elementToBeRemoved)

                self.rebalance_heaps()
        return result


def main():
    describe()

    slidingWindowMedian = SlidingWindowMedian()
    input = [1, 2, -1, 3, 5]
    print("Input : " + str(input) + " , 2")
    result = slidingWindowMedian.find_sliding_window_median(
        input, 2)
    print("Sliding window medians are: " + str(result))
    print("--------------")

    slidingWindowMedian = SlidingWindowMedian()
    print("Input : " + str(input) + " , 3")
    result = slidingWindowMedian.find_sliding_window_median(
        input, 3)
    print("Sliding window medians are: " + str(result))

main()
