###############################################################################
############################## LINKED LIST NODE ###############################
###############################################################################
class LinkedListNode
  attr_accessor :value, :next

  def initialize(value, nextPtr = nil)
    @value = value
    @next = nextPtr
  end

  def to_s
    res = self.value
    if block_given?
      res = yield self.value
    end
    res.to_s
  end
end


###############################################################################
################################# LINKED LIST #################################
###############################################################################
class LinkedList
  attr_accessor :head, :tail

  def initialize
    @head = nil
    @tail = nil
  end

  # Create new node with value and next as the current head
  #   Change head to be the new node
  #   If there is no tail, it means this is the first node - so make tail as the current node
  def prepend(val)
    node = LinkedListNode.new(val, self.head)
    self.head = node
    if !self.tail
      self.tail = node
    end
    self
  end

  # Create new node with value and next as nil
  #   Change current tail's next to be the new node
  #   Change current tail to become the new node
  #   If there is no head, it means this is the first node - so make head as the current node
  def append(val)
    node = LinkedListNode.new(val, nil)
    if !self.head
      self.head = node
    end
    if self.tail
      self.tail.next = node
    end
    self.tail = node
    self
  end

  # Start from the head and keep going in a loop to find the value
  def find(val)
    found = false
    if !self.head
      return nil
    end
    current = self.head
    while current.next
      if current.value == val
        found = true
        break
      end
      current = current.next
    end
    found ? current : nil
  end

  def delete(val)
    if !self.head
      return nil
    end
    deleted = nil
    previous = nil
    current = self.head
    while current.next
      if current.value == val
        deleted = current
        if current.next && previous # case when found node is between head and tail
          previous.next = current.next
        elsif !current.next && previous # case when found node is currently the tail
          self.tail = previous
        elsif current.next && !previous # case when found node is currently the head
          self.head = current.next
        else # case when found node is the only node in the linked list
          self.head = nil
          self.tail = nil
        end
        break
      end
      previous = current # Maintain a history variable to store the current node
      current = current.next # Then change current to the next node
    end
    deleted
  end

  # Store current head as result
  # Change head to next node of current head
  # Return result from above
  def delete_head
    if !self.head
      return nil
    end
    res = self.head
    if self.head.next
      self.head = self.head.next
    else
      self.head = nil
      self.tail = nil
    end
    res.next = nil
    res
  end

  # Traverse to the node before the current tail
  # Store the tail as the node to return
  # Change the tail to the node before the current tail and return the result from above
  def delete_tail
    res = self.tail
    # Handle case when there is only 1 node
    if self.head == self.tail
      self.head = nil
      self.tail = nil
      return res
    end
    # Traverse to the node before the current tail
    current = self.head
    while current.next.next
      current = current.next
    end
    # Make the last but one node the new tail
    current.next = nil
    self.tail = current
    res
  end

  def reverse
    return self unless self.head
    previous = nil
    current = self.head
    while current
      next_node = current.next # Save the next node of current to move forward in the loop
      current.next = previous # Set the next of the current node to the previous node
      previous = current # Update the history variable previous to be the current node
      current = next_node # Update the current node to next node saved from above
    end
    self.tail = self.head
    self.head = previous # The last previous is the last element of the iteration. So, make that the new head.
    self
  end

  def to_a
    nodes = []
    current = self.head
    while current
      nodes << current
      current = current.next
    end
    nodes
  end

  def empty?
    !self.head
  end

  def to_s
    self.to_a.collect {|node| node.to_s}
  end

  # construct from an array
  def self.from_array(values)
    ll = self.new
    values.each {|val| ll.append(val)}
    ll
  end
end


###############################################################################
################################## STACK ######################################
###############################################################################
class Stack
  attr_accessor :linked_list

  def initialize
    @linked_list = LinkedList.new
  end

  # construct from an array
  def self.from_array(values)
    s = self.new
    s.linked_list = LinkedList.from_array(values)
    s
  end

  def peek
    if self.linked_list.tail
      self.linked_list.tail.value
    else
      nil
    end
  end

  def push(val)
    self.linked_list.append(val)
    self
  end

  def pop
    res = self.linked_list.delete_tail
    return res ? res.value : nil
  end

  def empty?
    self.linked_list.empty?
  end

  def to_a
    self.linked_list.to_a.collect {|node| node.value}
  end

  def to_s
    self.to_a.to_s
  end
end


###############################################################################
################################## QUEUE ######################################
###############################################################################
class Queue
  attr_accessor :linked_list

  def initialize
    @linked_list = LinkedList.new
  end

  # construct from an array
  def self.from_array(values)
    s = self.new
    s.linked_list = LinkedList.from_array(values)
    s
  end

  def peek
    if self.linked_list.head
      self.linked_list.head.value
    else
      nil
    end
  end

  def enqueue(val)
    self.linked_list.append(val)
    self
  end

  def dequeue
    res = self.linked_list.delete_head
    return res ? res.value : nil
  end

  def empty?
    self.linked_list.empty?
  end

  def to_a
    self.linked_list.to_a.collect {|node| node.value}
  end

  def to_s
    self.to_a.to_s
  end
end


=begin
s = Stack.from_array([2,4,6,7,8,19,21,23,28])
puts "stack : #{s}"
s.push(43)
s.push(53)
puts "stack : #{s}"
puts "peek : #{s.peek}, stack : #{s}, pop : #{s.pop} stack : #{s}"

q = Queue.from_array([2,4,6,7,8,19,21,23,28])
puts "queue : #{q}"
q.enqueue(43)
q.enqueue(53)
puts "queue : #{q}"
puts "peek : #{q.peek}, queue : #{q}, dequeue : #{q.dequeue} queue : #{q}"
=end
