#  Problem : Given a linked list create a function to reverse it in place without using any additional data structures.
#  Input : 5 -> 3 -> 6 -> 4 -> 7 -> null
#  Output : 7 -> 4 -> 6 -> 3 -> 5 -> null

class Node:
    def __init__(self, data, next = None):
        self.data = data
        self.next = next

class LinkedList:
    def __init__(self, head = None):
        self.head = head

#  If we were allowed to use additional data structure we could have simply travered the existing list,
#  take each element and added it at the beginning of the list. So, somewhere along the middle it would have looked like this
#  6 -> 3 -> 5 -> null
#  To change links in place, maintain 3 pointers - previous, current and next
#      nextNode = current.next
#      current.next = previous
#      previous = current
#      current = nextNode
#  Time complexity of this is O(n)
def reverseList(list):
    if not list.head:
        return list
    previous = None
    current = list.head
    while current:
        nextNode = current.next
        current.next = previous
        previous = current
        current = nextNode
    list.head = previous
