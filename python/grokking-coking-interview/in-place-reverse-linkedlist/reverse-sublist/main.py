from __future__ import print_function
class Node():

    """Docstring for Node. """

    def __init__(self, value, node=None):
        self.value = value
        self.next = node

    def print_list(self):
        head = self
        while head is not None:
            print(head.value, end='->')
            head = head.next
        print('null')


def describe():
    desc = """
Problem : Given the head of a LinkedList and two positions p and q, reverse the LinkedList from position p to q
Example :
    Input: 1->2->3->4->5->null
    Output: 1->4->3->2->5->null

-------------
    """
    print(desc)

#  Time complexity is O(n)
def reverse_sub_list(head, p, q):
    count = 1 # indexes p and q are 1 based, so count starts from 1
    start,start_prev= None, None # find the starting position from where to reverse and keep track of the node before that index
    node = head
    while node is not None and count < p:
        start_prev = node
        node = node.next
        count += 1
    start = node

    prev = None
    current = start
    while current is not None and count <= q:
        next_node = current.next
        current.next = prev
        prev = current
        current = next_node
        count += 1
    stop = prev # This gives you the last node of the index q

    stop_next = current
    start_prev.next = stop
    start.next = stop_next

    return head

def main():
    describe()
    head = Node(1)
    head.next = Node(2)
    head.next.next = Node(3)
    head.next.next.next = Node(4)
    head.next.next.next.next = Node(5)

    print("Nodes of original LinkedList are: ", end='')
    head.print_list()
    result = reverse_sub_list(head, 2, 4)
    print("Nodes of reversed LinkedList are: ", end='')
    result.print_list()

main()
