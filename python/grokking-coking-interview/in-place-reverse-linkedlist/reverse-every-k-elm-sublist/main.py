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
Problem : Given the head of a Singly LinkedList, and a number k, reverse every k sized sub-list starting from the head.
          If, in the end, you are left with a sub-list with less than k elements, reverse that too.

Example :
    Input : 1->2->3->4->5->6->7->8->null
    Output: 3->2->1->6->5->4->8->7->null

-----------------

    """
    print(desc)

def reverse_every_k_elements(head, k):
    new_head = None
    sublist_start = None
    end_of_prev_sublist = None
    while head is not None:
        sublist_start = head

        # This section locally reverses the sublist
        # By the time the loop here exits, local_prev represents the head of the reversed sublist
        #  And the head represents the next node after the sublist ends
        local_prev = None
        count = 0
        while head is not None and count < k:
            node = head.next
            head.next = local_prev
            local_prev = head
            head = node
            count += 1

        sublist_start.next = head # Make the new sublist_start as the current position of head

        # If the end of previous sublist is null, then this is the first sublist reversal and will be the new head
        #  Otherwise the next node of the end of previous sublist will be the node representing the beginning of the reversed sublist
        if end_of_prev_sublist is None:
            new_head = local_prev
        else:
            end_of_prev_sublist.next = local_prev

        # Reposition the end of previous sublist
        end_of_prev_sublist = sublist_start

    return new_head



def alternate_impl_reverse_every_k_elements(head, k):
  if k <= 1 or head is None:
    return head

  current, previous = head, None
  while True:
    last_node_of_previous_part = previous
    # after reversing the LinkedList 'current' will become the last node of the sub-list
    last_node_of_sub_list = current
    next = None  # will be used to temporarily store the next node
    i = 0
    while current is not None and i < k:  # reverse 'k' nodes
      next = current.next
      current.next = previous
      previous = current
      current = next
      i += 1

    # connect with the previous part
    if last_node_of_previous_part is not None:
      last_node_of_previous_part.next = previous
    else:
      head = previous

    # connect with the next part
    last_node_of_sub_list.next = current

    if current is None:
      break
    previous = last_node_of_sub_list
  return head



def main():
    describe()
    head = Node(1)
    head.next = Node(2)
    head.next.next = Node(3)
    head.next.next.next = Node(4)
    head.next.next.next.next = Node(5)
    head.next.next.next.next.next = Node(6)
    head.next.next.next.next.next.next = Node(7)
    head.next.next.next.next.next.next.next = Node(8)

    print("Nodes of original LinkedList are: ", end='')
    head.print_list()
    result = reverse_every_k_elements(head, 3)
    print("Nodes of reversed LinkedList are: ", end='')
    result.print_list()

main()
