#  Problem : Sort a linked list in place

class Node:
    def __init__(self, data, next = None):
        self.data = data
        self.next = next

class LinkedList:
    def __init__(self, head = None):
        self.head = head


#  This is the bubble sort implementation
# The nested loop iomplementation here is a little cumbersome because we are changing links
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
