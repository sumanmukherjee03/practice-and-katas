class Node():

    """Node representing a node of a singly linked linkedlist"""

    def __init__(self, val, node = None):
        self.value = val
        self.next = None


def describe():
    desc = """
Problem : Given the head of a Singly LinkedList, write a method to return the middle node of the LinkedList.
If the total number of nodes in the LinkedList is even, return the second middle node.

Example : Input: 1 -> 2 -> 3 -> 4 -> 5 -> null
          Output: 3

          Input: 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> null
          Output: 4

-----------------
    """

#  Time complexity is O(n)
def find_middle_of_linked_list(head):
    slow, fast = head, head
    #  Based on this iteration, if linked list contains even number of Node, then fast will end up being None
    #  else if the number of Node is odd, fast will end up being the last node in the list
    #  Either way slow will always be in the middle because slow pointer travels at half the speed of fast
    while fast is not None and fast.next is not None:
        slow = slow.next
        fast = fast.next.next
    return slow

def main():
    describe()

    head = Node(1)
    head.next = Node(2)
    head.next.next = Node(3)
    head.next.next.next = Node(4)
    head.next.next.next.next = Node(5)

    print("Middle Node: " + str(find_middle_of_linked_list(head).value))

    head.next.next.next.next.next = Node(6)
    print("Middle Node: " + str(find_middle_of_linked_list(head).value))

main()
