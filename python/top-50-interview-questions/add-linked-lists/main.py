#  Problem - Given 2 linked lists representing 2 numbers in the reverse order, create
#  a function to return a linked list that adds the 2 numbers
#  For ex :
#       l1 = 3 -> 2 -> 1 -> null
#       l2 = 5 -> 9 -> 4 -> 3 -> null
#       return 8 -> 1 -> 6 -> 3 -> null
class Node:
    def __init__(self, data, next = None):
        self.data = data
        self.next = next

class LinkedList:
    def __init__(self, head = None):
        self.head = head

#  The time complexity is max(n, m) where n and m are the length of the 2 lists
def addTwoLinkedLists(list1, list2):
    res = LinkedList()
    ptr1 = list1.head
    ptr2 = list2.head
    head = res.head
    carry = 0
    while ptr1 is not None or ptr2 is not None:
        if head is None:
            head = Node(0)
            res.head = head
        else:
            head.next = Node(0)
            head = head.next
        if ptr1 is not None and ptr2 is not None:
            sum = ptr1.data + ptr2.data + carry
            ptr1 = ptr1.next
            ptr2 = ptr2.next
        elif ptr1 is not None and ptr2 is None:
            sum = ptr1.data + carry
            ptr1 = ptr1.next
        elif ptr1 is None and ptr2 is not None:
            sum = ptr2.data + carry
            ptr2 = ptr2.next
        head.data = sum % 10
        carry = int(sum/10)

    if carry > 0:
        head.next = Node(carry)
    return res
