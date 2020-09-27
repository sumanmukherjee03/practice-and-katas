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
def add_two_numbers(l1, l2)
	l3 = ListNode.new
	head = l3
	while l1 || l2 do
		val = l3.val
		val += l1.val if l1
		val += l2.val if l2

		rem = val % 10
		quo = val / 10
		puts "#{rem}, #{quo}"

		l3.val = rem

		l1 = l1.next if l1
		l2 = l2.next if l2

		if l1 || l2 || quo > 0
			l3.next = ListNode.new(quo, nil)
			l3 = l3.next
		end
	end
	return head
end
