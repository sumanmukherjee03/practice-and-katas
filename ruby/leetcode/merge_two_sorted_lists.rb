# Merge two sorted linked lists and return it as a sorted list. The list should be made by splicing together the nodes of the first two lists.
# Input: l1 = [1,2,4], l2 = [1,3,4]
# Output: [1,1,2,3,4,4]
# Example 2:

# Input: l1 = [], l2 = []
# Output: []
# Example 3:

# Input: l1 = [], l2 = [0]
# Output: [0]

# Definition for singly-linked list.
# class ListNode
#     attr_accessor :val, :next
#     def initialize(val = 0, _next = nil)
#         @val = val
#         @next = _next
#     end
# end
# @param {ListNode} l1
# @param {ListNode} l2
# @return {ListNode}
def merge_two_lists(l1, l2)
  # Guard clauses to handle empty lists
  return l2 unless l1
  return l1 unless l2

  head = l1
  previous_elm = nil

  while l1 && l2 do
    # If l1 is small move l1 forward
    if l1.val <= l2.val
      l1_next = l1.next
      previous_elm = l1
      l1 = l1_next
    elsif l1.val > l2.val # If l2 is small move l2 forward
      l2_next = l2.next
      if head == l1
        head = l2
      else
        previous_elm.next = l2
      end
      l2.next = l1
      previous_elm = l2
      l2 = l2_next
    end
  end

  if l2 && l2.val >= previous_elm.val
    previous_elm.next = l2
  end

  return head
end
