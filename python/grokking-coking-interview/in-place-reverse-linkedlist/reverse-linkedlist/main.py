from __future__ import print_function

class Node():
    def __init__(self, value, node=None):
        self.value = value
        self.next = node

    def print_list(self):
        node = self
        while node is not None:
            print(node.value, end='->')
            node = node.next
        print('null')
        print()


def describe():
    desc = """
Problem : Given the head of a Singly LinkedList, reverse the LinkedList.

Example :
-----------------

    """
    print(desc)

def reverse(head):
    prev = None # Keep track of the previous node in the list. Start with a null.
    while head is not None:
        node = head.next # Remember the next node before changing the pointers
        head.next = prev # Now change the pointer of the current head's next node to the previous node
        prev = head # Update the previous to be the current node
        head = node # Update the current head to be the next node that you remembered
    head = prev
    return head


def main():
    describe()

    head = Node(2)
    head.next = Node(4)
    head.next.next = Node(6)
    head.next.next.next = Node(8)
    head.next.next.next.next = Node(10)

    print("Nodes of original LinkedList are: ", end='')
    head.print_list()
    result = reverse(head)
    print("Nodes of reversed LinkedList are: ", end='')
    result.print_list()

main()
