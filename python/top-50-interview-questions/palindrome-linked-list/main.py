#  Problem : Given a linkd list of integers check if it is a palindrome with constant space complexity

class Node:
    def __init__(self, data, next = None):
        self.data = data
        self.next = next

class LinkedList:
    def __init__(self, head = None):
        self.head = head


# Had space complexity not been a constraint this could have been done by copying the linked list
# to another linked list by pushing elements at the front one after the other.
# And then once done by comparing element by element.
# Time complexity is O(n^2)
def isPalindromeList01(list):
    def getLinkedListLength(ll):
        arr = []
        n = ll.head
        l = 0
        while n:
            arr.append(n.data)
            l += 1
            n = n.next
        print(arr)
        return l

    # Maintain 2 pointers one for left and one for right
    # Increment the left pointer and decrement the right pointer and compare data
    listLength = getLinkedListLength(list)
    nodeLeft = list.head
    i = 0
    # Comparing till the middle is fine because both pointers are moving towards each other
    while i <= (listLength-1)/2:
        node = nodeLeft
        # If the left pointer moves forward by 1 node the right would have also decremented by 1 node
        # Hence the 2*i . And the -1 is because if there are n nodes there are n-1 links.
        for j in range(listLength - 2*i - 1):
            node = node.next
        nodeRight = node
        print(nodeLeft.data,nodeRight.data)
        if nodeLeft.data != nodeRight.data:
            return False
        nodeLeft = nodeLeft.next
        i += 1
    return True




# Use the slow and fast pointer technique here to find the middle of the linked list.
#  The slow pointer moves by one node while the fast pointer moves by 2 nodes, ie 2x the speed.
#  So, by the time the fast pointer reaches the end, the slow pointer has reached the middle.
#  Then from the middle node onwards, reverse the right half of the linked list.
#  Once the right half is reversed, compare node data from start of left half, ie the head of the list
#  with the start of the reversed right half of the list.
#  Time complexity of this solution is O(n) and no extra space is consumed other than vars, so space complexity of O(1)
def isPalindromeList02(list):
    def getMiddleNode(ll):
        slowPtr = ll.head
        fastPtr = ll.head
        while fastPtr and fastPtr.next:
            slowPtr = slowPtr.next
            fastPtr = fastPtr.next.next
        return slowPtr

    def reverseList(n):
        previous = None
        current = n
        while current:
            nextNode = current.next
            current.next = previous
            previous = current
            current = nextNode
        return previous


    mid = getMiddleNode(list)
    head = list.head
    revListHead = reverseList(mid)
    while revListHead:
        if head.data != revListHead.data:
            return False
        head = head.next
        revListHead = revListHead.next
    return True
