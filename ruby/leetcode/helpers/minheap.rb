# In this case a heap will be a binary min-heap,
# which means each parent node will be less than the value of each of, at most, two child nodes.
# The above structure can be represented in an array: [nil, 1, 3, 5, 8, 4]
# To find the index of the left child of a index i of parent node in the array use 2 * i.
# The index of the right child in the array is (2 * i) + 1.
# To find the index of parent for a child of index j in the array is j / 2.
# The above math works out if the root of the heap starts in the second container of the array, i.e., i = 1.

# MinHeap == PriorityQueue

# In heapsort the latest value is inserted in the next open leaf node going left to right.
# In the array representation this value would simply be inserted at the end.
# Now the heap has to be reordered so each parent is less than its children.
# This is accomplished be comparing the last inserted value to its parent.
# If it is less than its parent then those values are swapped
# This is done recursively up the tree until the it reaches the root or a parent with a value that is less than it.

# The root node of the heap represents the minimum value in the queue.
# Once this value is removed from the queue the next minimum value must occupy the root position.
# The root node and the last right-most node are swapped.
# Then the last right-most node is removed.
# Then the heap is reordered by comparing the root note to its smallest child and swapping them if the root is greater than that child.
# The the newly swapped child from the last operation is then compared to its smallest child and then again swapped if it is greater than that child.
# This is done recursively downward until the bottom-most node is reached or the child is greater than the parent.

class MinHeap
  # Implement a MinHeap using an array
  def initialize
    # Initialize arr with nil as first element
    # This dummy element makes it so that first real element is at index 1
    # You can now divide i / 2 to find an element's parent
    # Elements' children are at i * 2 && (i * 2) + 1
    @elements = [nil]
  end

  def min
    @elements[1..-1].min
  end

  def <<(element)
    # push item onto end (bottom) of heap
    @elements.push(element)
    # then bubble it up until it's in the right spot
    bubble_up(@elements.length - 1)
  end

  def pop
    # swap the first and last elements
    @elements[1], @elements[@elements.length - 1] = @elements[@elements.length - 1], @elements[1]
    # Min element is now at end of arr (bottom of heap)
    min = @elements.pop
    # Now bubble the top element (previously the bottom element) down until it's in the correct spot
    bubble_down(1)
    # return the min element
    min
  end

  def length
    @elements.length - 1
  end

  def peek
    @elements[1]
  end

  def print
    @elements
  end

  private

  def bubble_up(index)
    parent_i = index / 2
    # Done if you reach top element or parent is already smaller than child
    return if index <= 1 || @elements[parent_i] <= @elements[index]

    # Otherwise, swap parent & child, then continue bubbling
    @elements[parent_i], @elements[index] = @elements[index], @elements[parent_i]

    bubble_up(parent_i)
  end

  def bubble_down(index)
    child_i = index * 2
    return if child_i > @elements.size - 1

    # get the index of the largest child
    not_last = child_i < @elements.size - 1
    left = @elements[child_i]
    right = @elements[child_i + 1]
    child_i += 1 if not_last && right < left

    # stop if parent element is already less than children
    return if @elements[index] <= @elements[child_i]

    # otherwise, swap and continue
    @elements[index], @elements[child_i] = @elements[child_i], @elements[index]
    bubble_down(child_i)
  end
end
