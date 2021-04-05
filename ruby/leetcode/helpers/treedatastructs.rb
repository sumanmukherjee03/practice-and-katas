###############################################################################
############################## BINARY TREE NODE ###############################
###############################################################################
class BinaryTreeNode
  attr_accessor :left, :right, :parent, :value

  def initialize(value = nil)
    @left = nil
    @right = nil
    @parent = nil
    @value = value
  end

  # Recursively find height of left subtree + 1
  def left_height
    return 0 unless self.left
    self.left.height + 1
  end

  # Recursively find height of right subtree + 1
  def right_height
    return 0 unless self.right
    self.right.height + 1
  end

  # Find max of left height and right height
  def height
    [self.left_height, self.right_height].max
  end

  # Find diff of left height and right height
  def balance
    self.left_height - self.right_height
  end

  # Find the sum of all elements in the tree starting from this node
  def sum
    sum = self.value
    sum += self.left.sum if self.left
    sum += self.right.sum if self.right
    sum
  end

  # Reverse the tree as if the new version of the tree is a mirror object of itself
  def reverse
    l_reverse = self.left ? self.left.reverse : nil
    r_reverse = self.right ? self.right.reverse : nil
    self.set_left(r_reverse)
    self.set_right(l_reverse)
    self
  end

  # Only sets the value of a node, not it's pointers
  def set_value(val)
    self.value = val
    self
  end

  # Change the left subtree
  def set_left(node)
    if self.left
      self.left.parent = nil
    end
    self.left = node
    self.left.parent = self if node
    self
  end

  # Change the right subtree
  def set_right(node)
    if self.right
      self.right.parent = nil
    end
    self.right = node
    self.right.parent = self if node
    self
  end

  # Remove the child subtree from the matching node onwards
  def remove_child_subtree(node)
    if self.left && self.left.value == node.value
      self.left = nil
      return true
    end
    if self.right && self.right.value == node.value
      self.right = nil
      return true
    end
    return false
  end

  # Replace child subtree from matching node onwards with a new subtree containing new node
  def replace_child(target_node, replacement_node)
    return false unless target_node && replacement_node
    if self.left && self.left.value == target_node.value
      self.left = replacement_node
      return true
    end
    if self.right && self.right.value == target_node.value
      self.right = replacement_node
      return true
    end
    return false
  end

  # Inorder traversal Left, Right, Root
  def inorder
    traversal = []
    traversal.concat(self.left.inorder) if self.left
    traversal << self.value
    traversal.concat(self.right.inorder) if self.right
    traversal
  end

  # Preorder traversal Root, Left, Right
  def preorder
    traversal = []
    traversal << self.value
    traversal.concat(self.left.preorder) if self.left
    traversal.concat(self.right.preorder) if self.right
    traversal
  end

  # Postorder traversal Left, Right, Root
  def postorder
    traversal = []
    traversal.concat(self.left.postorder) if self.left
    traversal.concat(self.right.postorder) if self.right
    traversal << self.value
    traversal
  end

  # Levelorder traversal is essentially a depth first traversal of a tree
  def levelorder
    queue = [] # FIFO data structure
    traversal = []
    queue.push(self)
    while !queue.empty?
      node = queue.shift
      queue.push(node.left) if node.left
      queue.push(node.right) if node.right
      traversal.push(node.value)
    end
    traversal
  end

  # BFS or level order traversal but taking into account empty nodes so that it can be used for pretty printing
  def bfs_with_empty_nodes
    queue = [] # FIFO data structure
    traversal = []

    queue.push(self)
    while !queue.all? {|e| e == 'nil'}
      node = queue.shift
      if node != 'nil'
        if node.left
          queue.push(node.left)
        else
          queue.push('nil')
        end
        if node.right
          queue.push(node.right)
        else
          queue.push('nil')
        end
        traversal.push(node.value)
      else
        queue.push('nil')
        queue.push('nil')
        traversal.push('nil')
      end
    end
    traversal
  end

  def add_depth_to_output
    each_level = []
    res = []
    level = 0
    array = self.bfs_with_empty_nodes
    array.each do |elm|
      each_level.push(elm)
      if each_level.length > (2**level-1)
        res.push(each_level)
        each_level = []
        level += 1
      end
    end

    if each_level.length != 0
      until each_level.length >= 2**level
        each_level.push("nil")
      end
      res.push(each_level)
      each_level = []
    end

    res
  end

  def pretty_print
    array = self.add_depth_to_output
    if(array.length > 0)
      width=3 #max spacing of each node
      maxSpace = array[-1].length*width #total space in the last level
      array.each_with_index do |e,i|
        startSpace = maxSpace/(2**(i+1))
        space = maxSpace/(2**i)
        e.each_with_index do|k, j|
          if(k == "nil")
            k =" "
          end
          if(j==0)
            printf("%*s",startSpace,k)
          else
            printf("%*s",space,k)
          end
        end
        printf("\n")
      end
    end
  end

  # Sibling of parent
  def parents_sibling
    return nil unless self.parent
    return nil unless self.parent.parent
    return nil unless self.parent.parent.left || self.parent.parent.right
    if self.parent.value == self.parent.parent.left.value
      return self.parent.parent.right
    end
    self.parent.parent.left
  end

  def to_s
    self.value.to_s
  end

  # Class method to copy a node
  def self.copy(src, dest)
    dest.set_value(src.value)
    dest.set_left(src.left)
    dest.set_right(src.right)
  end

  # Algorithm to construct tree given it's inorder and preorder traversals
  #   1) Pick an element from preorder.
  #       - preorder always has the root of a tree/subtree first.
  #       - increment a preorder index variable to pick next element in next recursive call.
  #   2) create a new tree node with the data as picked element.
  #   3) find the picked element’s index in inorder
  #       - elements left of that index in inorder represent the left subtree nodes of the node with that value
  #       - elements right of that index in inorder represent the right subtree nodes of the node with that value
  #   4) recursively build tree for elements before the inorder index found above and make the built tree as left subtree of the node
  #   5) recursively build tree for elements after the inorder index and make the built tree as right subtree of the node
  #   6) return the node
  def self.from_arrays(inorder_arr, preorder_arr)
    klass = self
    preorder_index = 0

    # Dont use Proc.new - use lambda
    build_from_array = lambda do |in_arr|
      return nil if preorder_index >= preorder_arr.length # guard clause to stop
      val = preorder_arr[preorder_index] # get value for node
      preorder_index += 1 # increment preorder index to find the next root in the preorder traversal
      node = klass.new(val) # build node
      return node if in_arr.length == 1 # guard caluse to handle tree with just 1 node
      in_index = in_arr.index(val) # find index of the element from the root above, in inorder array
      ltree_vals = in_arr[0...in_index] # find elements in the left subtree
      rtree_vals = in_arr[in_index+1..-1] # find elements in the right subtree

      # build left subtree with elements from inorder array for the left subtree
      l_node = build_from_array.call(ltree_vals)
      node.set_left(l_node) if l_node

      # build right subtree with elements from inorder array for the right subtree
      r_node = build_from_array.call(rtree_vals)
      node.set_right(r_node) if r_node

      return node
    end

    return build_from_array.call(inorder_arr)
  end
end


###############################################################################
########################### BINARY SEARCH TREE NODE ###########################
###############################################################################
class BinarySearchTreeNode < BinaryTreeNode
  def insert(val)
    return self unless val
    # If value is lower than current nodes value go left
    #   Recursively try to insert the same value on the left subtree if a left subtree exists
    #   Otherwise create a new node and add it as the left child of current node
    if val < self.value
      return self.left.insert(val) if self.left
      node = BinarySearchTreeNode.new
      self.set_left(node)
      return self
    end
    # If value is greater than current nodes value go right
    #   Recursively try to insert the same value on the right subtree if a right subtree exists
    #   Otherwise create a new node and add it as the right child of current node
    if val > self.value
      return self.right.insert(val) if self.right
      node = BinarySearchTreeNode.new
      self.set_right(node)
      return self
    end
    self
  end

  # Find node matching value
  #   If current node has the asked value return that node
  #     Else recursively search left if value if lower than current nodes value
  #     Or recursively search right if value is greater than current nodes value
  def find(val)
    return self if self.value == val
    return self.left.find(val) if val < self.value && self.left
    return self.right.find(val) if val > self.value && self.right
    return nil
  end

  # Check if node or it's subtree contains a value
  def contains?(val)
    !!self.find(val)
  end

  def delete(val)
    node = self.find(val) # First, find the node to be deleted
    return false unless node # Guard clause - If the node to be deleted couldnt be found, return

    if !node.left && !node.right
      # If no left or rigth subtree
      #   If parent is there, remove this child from the parent
      #   Else this is the only node in the tree
      #     set the value of this node to nil
      if node.parent
        node.parent.remove_child_subtree(node)
      else
        node.set_value(nil)
      end

    elsif node.left && !node.right
      # Replace the current node with the node from the left subtree if it only has a left subtree
      node.parent.replace_child(node, node.left)

    elsif !node.left && node.right
      # Replace the current node with the node from the right subtree if it only has a right subtree
      node.parent.replace_child(node, node.right)

    else
      # When there are children in the right subtree
      #  - find the inorder successor of the node.
      #  - copy contents of the inorder successor to the node
      #  - delete the inorder successor node so that it's children can become direct children of the inorder successors parent
      # Inorder successor can be obtained by finding the minimum value in the right child of the node
      #
      # If there is a left and a right subtree
      #   If the next highest value is not the right child, ie the right child has a left subtree
      #     Then find the min value node from the right node onwards
      #     And then make the node to be deleted the node with the min value from above search
      #       ie, recursively delete the next big node
      #       and make the value of the current node as the next big node
      #   If the next highest value node is the right child
      #     Then make this node the right node
      #     And make the right subtree of this node the right subtree of the right child
      next_big_node = node.right.find_min
      if next_big_node.value != node.right.value
        self.delete(next_big_node.value)
        node.set_value(next_big_node.value)
      else
        node.set_value(node.right.value)
        node.set_right(node.right.right)
      end
    end
    return true
  end

  # Recursively traverse the left subtree to find the minimum value
  def find_min
    return self unless self.left
    self.left.find_min
  end

  # Recursively traverse the right subtree to find the max value
  def find_max
    return self unless self.right
    self.right.find_max
  end
end


###############################################################################
############################## BINARY SEARCH TREE #############################
###############################################################################
class BinarySearchTree
  attr_accessor :root

  def initialize
    @root = BinarySearchTreeNode.new
  end

  def insert(val)
    self.root.insert(val)
  end

  def delete(val)
    self.root.delete(val)
  end

  def contains?(val)
    self.root.contains?(val)
  end

  def inorder
    self.root.inorder
  end

  def preorder
    self.root.preorder
  end

  def postorder
    self.root.postorder
  end

  def find_min
    self.root.find_min
  end

  def find_max
    self.root.find_max
  end

  def height
    self.root.height
  end

  def balance
    self.root.balance
  end

  def to_a
    self.inorder
  end

  def to_s
    self.to_a.map {|node| node.to_s}
  end
end



###############################################################################
############################## SAMPLE RUN #############################
###############################################################################

=begin
n = BinaryTreeNode.from_arrays(["D","B","E","A","F","C"], ["A","B","D","E","C","F"])
puts "Tree : #{n.inorder}"
puts "Level order traversal : #{n.levelorder}"
n.pretty_print
n.reverse
puts "Inorder traversal of reverse tree : #{n.inorder}"
puts "Level order traversal of reverse tree : #{n.levelorder}"
n.pretty_print
puts "==========================="
m = BinaryTreeNode.from_arrays([8,4,10,9,11,2,5,1,6,3,7], [1,2,4,8,9,10,11,5,3,6,7])
puts "Tree : #{m.inorder}"
puts "Inorder traversal : #{m.inorder}"
puts "Preorder traversal : #{m.preorder}"
puts "Postorder traversal : #{m.postorder}"
puts "Level order traversal : #{m.levelorder}"
puts "BFS with empty nodes : #{m.bfs_with_empty_nodes}"
puts "Sum : #{m.sum}"
m.pretty_print
puts "------------------------------"
m.reverse
puts "Level order traversal of reverse tree : #{m.levelorder}"
m.pretty_print
=end
