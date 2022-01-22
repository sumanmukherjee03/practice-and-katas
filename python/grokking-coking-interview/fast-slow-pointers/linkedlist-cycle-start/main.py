class Node:
    def __init__(self, val, node=None):
        self.value = val
        self.next = node

    def print_list(self):
        node = self
        while node is not None:
            print(node.value, "->")
            node = node.next
        print()



def describe():
    desc = """
Problem : Given the head of a Singly LinkedList that contains a cycle, write a function to find the starting node of the cycle.
Example :
    Input : 1 -> 2 -> 3 -> 4 -> 5 -> 6
                    ^              |
                    |--------------|
    Output : Linked List cycle starts at 3
------------------

    """
    print(desc)

#  Time complexity is O(n)
def find_cycle_start(head):
    slow, fast = head, head
    cycle_len = 0
    while fast is not None and fast.next is not None:
        slow = slow.next
        fast = fast.next.next
        if slow == fast:
            ptr = slow.next
            cycle_len = 1 # Start the count from 1 because we need to include the node of the slow pointer
            while ptr != slow:
                cycle_len += 1
                ptr = ptr.next
            break

    if cycle_len == 0:
        return None

    temp = head
    for i in range(0, cycle_len):
        temp = temp.next

    #  Now start 2 pointers - one from the beginning and one from cycle length nodes ahead and increment them 1 at a time until they meet
    ptr1, ptr2 = head, temp
    while ptr1 != ptr2:
        ptr1 = ptr1.next
        ptr2 = ptr2.next

    return ptr1


def main():
    describe()
    head = Node(1)
    head.next = Node(2)
    head.next.next = Node(3)
    head.next.next.next = Node(4)
    head.next.next.next.next = Node(5)
    head.next.next.next.next.next = Node(6)
    head.print_list()

    head.next.next.next.next.next.next = head.next.next
    print("linkedlist cycle start: " + str(find_cycle_start(head).value))

    head.next.next.next.next.next.next = head.next.next.next
    print("linkedlist cycle start: " + str(find_cycle_start(head).value))

    head.next.next.next.next.next.next = head
    print("LinkedList cycle start: " + str(find_cycle_start(head).value))

main()
