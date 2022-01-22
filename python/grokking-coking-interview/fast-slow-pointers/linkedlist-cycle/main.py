class Node:
    def __init__(self, value, next=None):
        self.value = value
        self.next = next


def describe():
    desc = """
Problem : Given the head of a Singly LinkedList, write a function to determine if the LinkedList has a cycle in it or not
Example :
    Input : 1 -> 2 -> 3 -> 4 -> 5 -> 6
                    ^              |
                    |--------------|
    Output : Linked List has cycle
------------------

    """
    print(desc)

#  Time complexity is O(n)
def has_cycle(head):
    slow = head
    fast = head
    #  As long there is a fast pointer pointing to non-null node there will always be a slow pointer pointing to non-null node
    #  If fast.next is null, then you have hot the end of the linked list already
    while fast is not None and fast.next is not None:
        slow = slow.next # move by 1 node
        fast = fast.next.next # move by 2 nodes
        if slow == fast:
            return True
    return False

#  Time complexity O(n)
#  To find the cycle length, first find the cycle. That means the slow and fast pointers must meet.
#  Once they meet keep the slow pointer intact and use another pointer that proceeds one node at a time
#  and goes through the cycle and meets the slow pointer again. That will give us the length of the cycle.
def find_cycle_length(head):
    slow, fast = head, head
    while fast is not None and fast.next is not None:
        slow = slow.next
        fast = fast.next.next
        if slow == fast:
            return calculate_cycle_length(slow)
    return 0

def calculate_cycle_length(slow):
    ptr = slow.next
    count = 1
    while ptr != slow:
        ptr = ptr.next
        count += 1
    return count





def alternate_implementation_calculate_cycle_length(slow):
  current = slow
  cycle_length = 0
  while True:
    current = current.next
    cycle_length += 1
    if current == slow:
      break
  return cycle_length




def main():
    describe()
    head = Node(1)
    head.next = Node(2)
    head.next.next = Node(3)
    head.next.next.next = Node(4)
    head.next.next.next.next = Node(5)
    head.next.next.next.next.next = Node(6)
    print("LinkedList has cycle: " + str(has_cycle(head)))

    head.next.next.next.next.next.next = head.next.next
    print("LinkedList has cycle: " + str(has_cycle(head)))

    head.next.next.next.next.next.next = head.next.next.next
    print("LinkedList has cycle: " + str(has_cycle(head)))
    print("LinkedList has cycle of length: " + str(find_cycle_length(head)))

main()
