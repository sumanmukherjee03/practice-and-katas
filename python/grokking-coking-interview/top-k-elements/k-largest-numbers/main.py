import heapq

def describe():
    desc = """
Problem : Given an unsorted array find the K largest numbers in it.
Example :
    Input : [3, 1, 5, 12, 2, 11], K = 3
    Output : [5, 12, 11]


---------------
    """
    print(desc)

#  The bruteforce solution is to sort the given array and return the k largest numbers from the end.
#  A better solution is to iterate from the array and keep inserting the numbers into a max heap and then pop k times from the heap
def unoptimized_find_k_largest_numbers(nums, k):
    res = []
    maxheap = []
    for i in range(0, len(nums)):
        heapq.heappush(maxheap, -nums[i])
    for j in range(0, k):
        res.append(-heapq.heappop(maxheap))
    return res

#  Another more optimized solution is to maintain a min heap of k elements as we iterate through the given array.
#  As we encounter numbers, we take out the smallest element from the maxheap and insert the larger of the 2 numbers,
#  ie the smallest from the heap and the new number in the iteration into the heap.
#  That way at the end of the iteration even though it is a minheap it will contain the largest elements of the array.
#  Time complexity for this is O(k log(k) + (n-k) log(k))
def find_k_largest_numbers(nums, k):
    minheap = []
    for i in range(0, len(nums)):
        if len(minheap) < k:
            heapq.heappush(minheap, nums[i])
        else:
            if nums[i] > minheap[0]:
                heapq.heappop(minheap)
                heapq.heappush(minheap, nums[i])
    return minheap

def main():
    describe()
    input, key = [3, 1, 5, 12, 2, 11], 3
    print("Input : " + str(input) + " , " + str(key))
    print("Output : " + str(find_k_largest_numbers(input, key)))

    input, key = [5, 12, 11, -1, 12], 3
    print("Input : " + str(input) + " , " + str(key))
    print("Output : " + str(find_k_largest_numbers(input, key)))

main()
