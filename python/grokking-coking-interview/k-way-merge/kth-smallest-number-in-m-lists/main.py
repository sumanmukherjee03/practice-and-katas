import heapq

def describe():
    desc= """
Problem :
    Given m sorted arrays find the kth smallest number among all arrays
Similar problems :
    Given m sorted arrays, find the median number among all arrays. This is the same problem with k = n/2.
    Given a list of m sorted arrays, merge them into one sorted list.
Example :
    Input: L1=[2, 6, 8], L2=[3, 6, 7], L3=[1, 3, 4], K=5
    Output: 4
    Explanation: The 5th smallest number among all the arrays is 4, this can be verified from the merged list of all the arrays: [1, 2, 3, 3, 4, 6, 6, 7, 8]

--------------
    """
    print(desc)


#  To find the kth smallest number among m sorted arrays, if you maintain a maxheap of size k
#  which contains all the smallest numbers, then the top element of the heap will me the largest of the smallest numbers.
#  Go through all the elements of all the lists and keep a maxheap of the smallest numbers and finally pop from that heap to get the result.
def unoptimized_find_kth_smallest_number(lists, k):
    maxheap = []
    m = len(lists)
    for i in range(0, m):
        l = lists[i]
        for j in range(0, len(l)):
            if len(maxheap) < k:
                heapq.heappush(maxheap, -l[j])
            else:
                if -maxheap[0] > l[j]:
                    heapq.heappop(maxheap)
                    heapq.heappush(maxheap, -l[j])
    return -heapq.heappop(maxheap)




#  This solution is similar to the problem for the merging k sorted linked lists.
#  So, for this, we maintain a minheap of size equal to the length of the number of lists.
#  We keep popping numbers from the minheap and enter them into the resulting merged list.
#  And we perform that operation until we reach the count of k elements. Then that element is the kth smallest element.
#  The only difference is that we need to keep track of which list and what index in that list was popped from the minheap.
#  Time complexity of this solution is O(k * log(m))
def find_kth_smallest_number(lists, k):
    minheap = []
    for i in range(0, len(lists)):
        heapq.heappush(minheap, (lists[i][0], 0, lists[i]))

    count, number = 0, 0
    while minheap:
        number, j, l = heapq.heappop(minheap)
        count += 1
        if count == k:
            break
        if len(l) > j+1:
            heapq.heappush(minheap, (l[j+1], j+1, l))
    return number

def main():
    describe()
    input = [[2, 6, 8], [3, 6, 7], [1, 3, 4]]
    k = 5
    print("Input : " + str(input) + " , " + str(k))
    print("Kth smallest number is: " + str(find_kth_smallest_number(input, k)))

main()
