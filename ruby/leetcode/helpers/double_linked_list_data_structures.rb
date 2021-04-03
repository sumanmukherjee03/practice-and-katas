###############################################################################
########################### DOUBLE LINKED LIST NODE ###########################
###############################################################################
class DoubleLinkListNode
  require 'securerandom'

  attr_accessor :id, :value, :prev, :next

  def initialize(value, prevPtr = nil, nextPtr = nil)
    @id = SecureRandom.hex
    @value = value
    @prev = prevPtr
    @next = nextPtr
  end

  def to_s
    res = self.value
    if block_given?
      res = yield self.value
    end
    res.to_s
  end

  def eql?(other)
    return other.id == self.id
  end
  alias :== :eql?
end

###############################################################################
######################### DOUBLE LINK LINKED LIST ###########################
###############################################################################
class DoubleLinkLinkedList
  attr_accessor :head, :tail

  def initialize
    @head = nil
    @tail = nil
  end

  def append(val)
    node = DoubleLinkListNode.new(val, self.tail, nil)
    unless self.head
      self.head = node
    end
    if self.tail
      self.tail.next = node
    end
    self.tail = node
    self
  end

  def prepend(val)
    node = DoubleLinkListNode.new(val, nil, self.head)
    unless self.tail
      self.tail = node
    end
    if self.head
      self.head.prev = node
    end
    self.head = node
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
    previous_node = nil
    current_node = self.head
    while current_node.next
      if current_node.value == val
        deleted = current_node
        if current_node.next && previous_node # case when found node is between head and tail
          next_node = current_node.next
          previous_node.next = next_node
          next_node.prev = previous_node
        elsif !current_node.next && previous_node # case when found node is currently the tail
          next_node = nil
          previous_node.next = nil
          self.tail = previous_node
        elsif current_node.next && !previous # case when found node is currently the head
          next_node = current_node.next
          next_node.prev = nil
          self.head = next_node
        else # case when found node is the only node in the linked list
          self.head = nil
          self.tail = nil
        end
        break
      end
      previous_node = current_node # Maintain a history variable to store the current node
      current_node = current_node.next # Then change current to the next node
    end
    deleted
  end

  def delete_head
    unless self.head
      return nil
    end
    res = self.head
    if self.head.next
      new_head = self.head.next
      new_head.prev = nil
      self.head = new_head
    else
      self.head = nil
      self.tail = nil
    end
    res.next = nil
    res
  end

  def delete_tail
    unless self.tail
      return nil
    end
    res = self.tail
    if self.tail.prev
      new_tail = self.tail.prev
      new_tail.next = nil
      self.tail = new_tail
    else
      self.tail = nil
      self.head = nil
    end
    res.prev = nil
    res
  end

  def reverse
    return self unless self.head
    current_node = self.head
    self.tail = current_node
    previous_node = nil
    while current_node.next
      next_node = current.next
      current_node.next = previous_node
      current_node.previous = next_node
      previous_node = current_node
      current_node = next_node
    end
    self.head = previous_node
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
    dll = self.new
    values.each {|val| dll.append(val)}
    dll
  end
end


###############################################################################
################################## DEQUE ######################################
###############################################################################
class Deque
  attr_accessor :double_link_linked_list

  def initialize
    @double_link_linked_list = DoubleLinkLinkedList.new
  end

  # construct from an array
  def self.from_array(values)
    dq = self.new
    dq.double_link_linked_list = DoubleLinkLinkedList.from_array(values)
    dq
  end

  def push(val)
    self.double_link_linked_list.append(val)
    self
  end
  alias :enqueue :push

  # Remove last element and return that
  def pop
    node = self.double_link_linked_list.delete_tail
    node ? node.value : nil
  end

  # Remove first element and return that
  def shift
    node = self.double_link_linked_list.delete_head
    node ? node.value : nil
  end

  # shift right by adding an element at the beginning
  def unshift(val)
    self.double_link_linked_list.prepend(val)
  end

  def first
    self.double_link_linked_list.head ? self.double_link_linked_list.head.value : nil
  end

  def last
    self.double_link_linked_list.tail ? self.double_link_linked_list.tail.value : nil
  end

  def empty?
    self.double_link_linked_list.empty?
  end

  def to_s
    self.double_link_linked_list.to_s
  end

  def to_a
    self.double_link_linked_list.to_a
  end
end
