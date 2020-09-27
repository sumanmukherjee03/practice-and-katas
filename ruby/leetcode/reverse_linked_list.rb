# Definition for singly-linked list.
# class ListNode
#     attr_accessor :val, :next
#     def initialize(val = 0, _next = nil)
#         @val = val
#         @next = _next
#     end
# end
# @param {ListNode} head
# @return {ListNode}
def reverse_list(head)
    new_head = nil
    while head do
        node = ListNode.new(head.val, nil)
        if !new_head
            new_head = node
        else
            node.next = new_head
            new_head = node
        end
        head = head.next
    end
    return new_head
end
