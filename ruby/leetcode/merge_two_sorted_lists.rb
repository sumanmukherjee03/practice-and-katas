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
    return l2 unless l1
    return l1 unless l2

    head = l1
    previous_l1 = nil

    while l1 && l2 do
        # If l1 is small move l1 forward
        if l1.val <= l2.val
            previous_l1 = l1
            l1 = l1.next
        elsif l1.val > l2.val # If l2 is small move l2 forward
            l2_next = l2.next
            if head == l1
                head = l2
            else
                previous_l1.next = l2
            end
            l2.next = l1
            previous_l1 = l2
            l2 = l2_next
        end
    end

    if l2 && l2.val >= previous_l1.val
        previous_l1.next = l2
    end

    return head
end
