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

main()
