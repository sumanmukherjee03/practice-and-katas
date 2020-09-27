# Definition for singly-linked list.
# class ListNode
#     attr_accessor :val, :next
#     def initialize(val = 0, _next = nil)
#         @val = val
#         @next = _next
#     end
# end
# @param {ListNode[]} lists
# @return {ListNode}
def merge_k_lists(lists)
  return nil if lists.empty?
  return lists.first if lists.length == 1

  two_way_merge = lambda do |l1,l2|
    return l1 if !l1 && !l2
    return l2 if !l1 && l2
    return l1 if l1 && !l2

    l = ListNode.new
    head = l
    while l1 || l2 do
      if (l1 && l2 && l1.val < l2.val) || (l1 && !l2)
        l.val = l1.val
        l1 = l1.next
      elsif (l1 && l2 && l1.val >= l2.val) || (!l1 && l2)
        l.val = l2.val
        l2 = l2.next
      end
      if l1 || l2
        l.next = ListNode.new
        l = l.next
      end
    end
    return head
  end

  i = 1
  res = lists[0]
  while i < lists.length do
    res = two_way_merge.call(res,lists[i])
    i += 1
  end

  return res
end
