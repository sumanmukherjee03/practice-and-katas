import heapq

class ListNode:
    def __init__(self, value):
        self.value = value
        self.next = None
        self._list_index = -1 #  We are maintaining this information here because we want to track which list this node belongs to

    #  Implement the less than operator so that we can insert these nodes in a minheap
    def __lt__(self, other):
        return self.value < other.value

    def set_list_index(self, idx):
        self._list_index = idx

    def get_list_index(self):
        return self._list_index



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
    result, head = None, None
    minheap = []
    for i in range(0, len(lists)):
        lists[i].set_list_index(i)
        heapq.heappush(minheap, lists[i]) #  Insert the list node into minheap
        lists[i] = lists[i].next #  Move the head of the head of the list forward

    #  I am keeping this lambda here just for reference later : any(map(lambda x : x.value is not None, lists))

    #  We will need to iterate until the minheap is empty
    while len(minheap) > 0:
        elm = heapq.heappop(minheap)
        list_index = elm.get_list_index()
        node = lists[list_index]

        if result is None:
            result = ListNode(elm.value)
            head = result
        else:
            head.next = ListNode(elm.value)
            head = head.next

        if node is not None:
            node.set_list_index(list_index)
            heapq.heappush(minheap, node)
            lists[list_index] = lists[list_index].next

    return result


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
