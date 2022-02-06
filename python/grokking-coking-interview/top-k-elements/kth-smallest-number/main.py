import heapq

def describe():
    desc = """
Problem: Given an unsorted array, find the kth smallest number in it.
Example :
    Input: [1, 5, 12, 2, 11, 5], K = 3
    Output: 5
    Explanation: The 3rd smallest number is 5, as the first two smaller numbers are [1, 2]

-----------------
    """
    print(desc)

#  If you maintain a minheap of size (len(nums) - (k - 1)) that contains the largest elements in the array,
#  then whats left in the array is the (k - 1) smaller elements and the top of the heap will contain the kth smallest element
def suboptimal_find_kth_smallest_number(nums, k):
    heap_len = len(nums) - (k-1)
    minheap = []
    for i in range(0, len(nums)):
        if len(minheap) < heap_len:
            heapq.heappush(minheap, nums[i])
        else:
            if nums[i] > minheap[0]:
                heapq.heappop(minheap)
                heapq.heappush(minheap, nums[i])
    return heapq.heappop(minheap)

#  Alternatively another approach is to maintain a maxheap of k elements.
#  Each time we iterate upon an element of the array we replace the top of the heap with the smaller of the top of heap and the current number.
#  That way the max will have the k smallest elements in the array
#  Complexity of this is O(n*log(k))
def find_kth_smallest_number(nums, k):
    maxheap = []
    for i in range(0, len(nums)):
        if len(maxheap) < k:
            heapq.heappush(maxheap, -nums[i])
        else:
            if nums[i] < -maxheap[0]:
                heapq.heappop(maxheap)
                heapq.heappush(maxheap, -nums[i])
    return -heapq.heappop(maxheap)

#  Another approach is to insert all elements into a minheap and then take out k elements from the minheap
#  the complexity of that approach would be O(n + k*log(n))
def main():
    describe()
    input, key = [1, 5, 12, 2, 11, 5], 3
    print("Input : " + str(input) + " , " + str(key))
    print("Kth smallest number is: " + str(find_kth_smallest_number(input, key)))

    input, key = [1, 5, 12, 2, 11, 5], 4
    print("Input : " + str(input) + " , " + str(key))
    print("Kth smallest number is: " + str(find_kth_smallest_number(input, key)))

    input, key = [5, 12, 11, -1, 12], 3
    print("Input : " + str(input) + " , " + str(key))
    print("Kth smallest number is: " + str(find_kth_smallest_number(input, key)))

main()
