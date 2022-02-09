import heapq

class ListNode:
    def __init__(self, value):
        self.value = value
        self.next = None

    #  Implement the less than operator so that we can insert these nodes in a minheap
    def __lt__(self, other):
        return self.value < other.value



def describe():
    desc = """
Problem : Given an array of k sorted linkedlists, merge them into one sorted list
Example :
    Input: L1=[2, 6, 8], L2=[3, 6, 7], L3=[1, 3, 4]
    Output: [1, 2, 3, 3, 4, 6, 6, 7, 8]

    Input: L1=[5, 8, 9], L2=[1, 7]
    Output: [1, 5, 7, 8, 9]

--------------
    """
    print(desc)


#  Remember these lists are already sorted.
#  So at any point we are only concerned with making a decision between k elements with each element belonging to a different list.
#  So maintain a minheap will k elements will do the job.
#  Pop from the minheap to get the min element and insert it into a new list.
#  Then for the element which was popped, get the list which it belonged to and move the head forward to the next node.
#  Then insert that node back into the minheap.
def merge_lists(lists):
    head, tail = None, None
    minheap = []
    for i in range(0, len(lists)):
        if lists[i] is not None:
            heapq.heappush(minheap, lists[i]) #  Insert the list node into minheap

    #  We will need to iterate until the minheap is empty
    #  Pop from the minheap and move the pointer for that node ahead to the next node
    #  Then insert the node that you popped into the new list.
    #  Retain the next pointer of the node you entered into the resulting list.
    #  Then insert that next node into the minheap.
    #  Remember to retain the next pointer of the node you are entering into the minheap.
    #  At no point mess with the next pointers because they help you move forward when you pop from the list and these nodes are all references.
    while len(minheap) > 0:
        node = heapq.heappop(minheap)
        if head is None:
            head = tail = node
        else:
            tail.next = node
            tail = tail.next

        if node.next is not None:
            heapq.heappush(minheap, node.next)

    return head

#  Total time complexity is O(n log(k)) where k is the number of lists and n is the total number of nodes in the list
def main():
    describe()
    print("Input : ")
    print("L1 : 2->6->8")
    l1 = ListNode(2)
    l1.next = ListNode(6)
    l1.next.next = ListNode(8)

    print("L2 : 3->6->7")
    l2 = ListNode(3)
    l2.next = ListNode(6)
    l2.next.next = ListNode(7)

    print("L3 : 1->3->4")
    l3 = ListNode(1)
    l3.next = ListNode(3)
    l3.next.next = ListNode(4)
    print("--------------------")

    result = merge_lists([l1, l2, l3])
    print("Here are the elements form the merged list: ", end='')
    while result != None:
        print(str(result.value) + " ", end='')
        result = result.next


main()

