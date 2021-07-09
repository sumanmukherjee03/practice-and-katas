#  Problem : Sort a linked list in place

class Node:
    def __init__(self, data, next = None):
        self.data = data
        self.next = next

class LinkedList:
    def __init__(self, head = None):
        self.head = head


# This is the bubble sort implementation
# The nested loop iomplementation here is a little cumbersome because we are changing links instead of swapping data
# Time complexity O(n^2)
def sortList01(list):
    if not list.head:
        return
    t = list.head
    while t:
        node = list.head
        previous = None
        while node.next:
            n1 = node.next
            n2 = node.next.next # node.next.next stays on the right side of the equal to, so even if it's None it doesnt matter
            # Swap the nodes if the current node and it's next node are unsorted
            if node.data > n2.data:
                if previous:
                    previous.next = n1
                n1.next = node
                node.next = n2
                previous = n1
            else:
                previous = node
                node = node.next
        t = t.next



# Here is the same bubble sort implementation as above but with swapping of data and not the actual nodes, ie not link exchanges.
# Time complexity O(n^2)
def sortList02(list):
    if not list.head:
        return
    i = list.head
    while i:
        j = list.head
        while j.next:
            # Swap the nodes if the current node and it's next node are unsorted
            if j.data > j.next.data:
                j.data, j.next.data = j.next.data, j.data
            j = j.next # Here increment the pointer and move to the next node because we havent changed the links here
        i = i.next


#  Time complexity of this is O(nlogn)
def sortList03(list):
    if not list.head:
        return list
    if not list.head.next:
        return list

    def getMidNode(node):
        slow = node
        fast = node
        while fast.next.next:
            slow = slow.next
            fast = fast.next.next
        return slow

    def mergeLists(l, r):
        res = LinkedList.new()
        resPtr = None
        prev_l = None
        prev_r = None
        while l and r:
            node = Node.new()
            if l.data < r.data:
                node.data = l.data
                prev_l = l
                l = l.next
            else:
                node.data = r.data
                prev_r = r
                r = r.next
            if not res.head:
                res.head = node
                resPtr = res.head
            else:
                resPtr.next = node
        if l:
            prev_l.next = None
            resPtr.next = l
        else:
            prev_r.next = None
            resPtr.next = r
        return res

    mid = getMidNode(list.head)
    left = list.head
    right = mid.next
    mid.next = None
    return mergeLists(sortList03(left), sortList03(right))
