# Merge two sorted linked lists and return it as a sorted list. The list should be made by splicing together the nodes of the first two lists.
# Input: l1 = [1,2,4], l2 = [1,3,4]
# Output: [1,1,2,3,4,4]
# Example 2:

# Input: l1 = [], l2 = []
# Output: []
# Example 3:

# Input: l1 = [], l2 = [0]
# Output: [0]


# We can recursively define the result of a merge operation on two lists.

# list1[0] + merge(list1[1:], list2) if list1[0] < list2[0]
# list2[0] + merge(list1,list2[1:]) otherwise

# Namely, the smaller of the two lists' heads plus the result of a merge on the rest of the elements.

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
  return l1 unless l2
  return l2 unless l1

  if l1.val < l2.val
    l1.next = merge_two_lists(l1.next, l2)
    return l1
  else
    l2.next = merge_two_lists(l1,l2.next)
    return l2
  end
end

# We can achieve the same idea via iteration

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

  result = ListNode.new()
  head = result

  while l1 && l2
    if l1.val < l2.val
      head.next = l1
      l1 = l1.next
    else
      head.next = l2
      l2 = l2.next
    end
    head = head.next
  end

  head.next = l1 if l1
  head.next = l2 if l2

  result = result.next

  return result
end
