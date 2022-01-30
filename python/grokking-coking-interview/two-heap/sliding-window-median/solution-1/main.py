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

    def insert_num(self, num):
        if not self.maxheap or -self.maxheap[0] >= num:
            heapq.heappush(self.maxheap, -num)
        else:
            heapq.heappush(self.minheap, num)
        #  Rebalnce the heaps after insertion. The max heap will alaways have either equal number of
        #  elements as the min heap or 1 element more than it
        #  Rebalancing after insertion and removal ensres that.
        self.rebalance_heaps()

    def remove(self, element):
        if element <= -self.maxheap[0]:
            self._remove_from_partition(-element, self.maxheap)
        else:
            self._remove_from_partition(element, self.minheap)
        self.rebalance_heaps()

    def _remove_from_partition(self, element, heap):
        idx = heap.index(element)
        #  Delete the element by exchanging it with the last element of the array
        #  And then deleting from the last element of the array
        #  However, this means we will be losing the heap property because the previously placed last element is now at a wrong position
        #  So we will need to re-heapify the data structures
        heap[idx] = heap[-1]
        del heap[-1]
        if idx < len(heap):
            heapq.heapify(heap)
            #  heapq._siftup(heap, idx)
            #  heapq._siftdown(heap, 0, idx)

    def find_median(self):
        median = 0
        if len(self.maxheap) == len(self.minheap):
            median = (-self.maxheap[0] + self.minheap[0])/2.0
        else:
            median = -self.maxheap[0]
        return median


    def rebalance_heaps(self):
        if len(self.minheap) > len(self.maxheap):
            elm = heapq.heappop(self.minheap)
            heapq.heappush(self.maxheap, -elm)
        elif len(self.maxheap) > len(self.minheap)+1:
            elm = -heapq.heappop(self.maxheap)
            heapq.heappush(self.minheap, elm)
        return

    def find_sliding_window_median(self, arr, k):
        res = [0.0 for x in range(len(arr) - k + 1)] # Initialize the resulting array here with 0

        #  Populate the min and max heaps for the first sliding window and calculate the result
        #  so that from next iteration on elements you can simply add and remove elements as the sliding window moves
        for i in range(0, k):
            self.insert_num(arr[i])
        res[0] = self.find_median()

        for i in range(1, len(arr)-k+1):
            #  Remove i-1 th number from maxheap or minheap, insert the new num
            #  The rebalancing is taken care of by the insertion and removal
            self.remove(arr[i-1])
            self.insert_num(arr[i+k-1])
            res[i] = self.find_median()
        return res


def main():
    describe()

    slidingWindowMedian = SlidingWindowMedian()
    input = [1, 2, -1, 3, 5]
    result = slidingWindowMedian.find_sliding_window_median(
        input, 2)
    print("Input : " + str(input))
    print("Sliding window medians are: " + str(result) + " , 2")
    print("--------------")

    slidingWindowMedian = SlidingWindowMedian()
    result = slidingWindowMedian.find_sliding_window_median(
        input, 3)
    print("Input : " + str(input) + " , 3")
    print("Sliding window medians are: " + str(result))

main()
