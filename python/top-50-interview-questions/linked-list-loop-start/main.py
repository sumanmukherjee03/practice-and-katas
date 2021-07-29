#  Problem : Given a linked list, create 2 functions.
#  First function should return if the linked list has a loop or not
#  Second function should return where the loop starts in the linked list
class Node:
    def __init__(self, data, next = None):
        self.data = data
        self.next = next

class LinkedList:
    def __init__(self, head = None):
        self.head = head

#  If there is a loop then at some point the hare and tortoise will meet,
#  because otherwise the hare would have gone much ahead of the tortoise and would have reached the end or end but 1
#  If distance from head of linked list to start of loop (inclusive) in linked list is D
#  and the distance between start of linked list to the node where the 2 pointers meet is K, not including the start of loop
#  and the number of iterations of the tortoise is N
#  Then :
    #  For the hare
        #  2N = D + K + j * C   - where j is an arbitrary number of cycles of length C in the loop
    #  For the tortoise
        #  N = D + K + i * C    - where i is an arbitrary number of cycles of length C in the loop
    #  So, if we subtract the 2 equations above
        #  N = (j-i) * C
#  Meaning, there exists a number N of our choosing which will equal to a few cycles
#  That justifies that the 2 pointers for hare and tortoise will meet at some point.
def findMeetingPointofPtrs(list):
    if not list.head:
        return False
    tortoise = list.head
    hare = list.head
    found = None
    while hare is not None:
        if hare == tortoise:
            found = tortoise
            break
        else:
            if hare.next:
                hare = hare.next.next
            else:
                hare = None
            tortoise = tortoise.next
    return found

def hasLoop(list):
    node = findMeetingPointofPtrs(list)
    if not node:
        return False
    else:
        return True

# A fairly good explanation exists in this video - https://www.youtube.com/watch?v=-YiQZi3mLq0
def startOfLoop(list):
    if list.head is None:
        return None
    ptr1 = list.head
    ptr2 = findMeetingPointofPtrs(list)
    if ptr2 is None:
        return None
    found = None
    while ptr1 != ptr2:
        if ptr1 == ptr2:
            found = ptr1
            break
        ptr1 = ptr1.next
        ptr2 = ptr2.next
    return found
