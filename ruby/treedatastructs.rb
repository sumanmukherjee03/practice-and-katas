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
    self.left.parent = self
    self
  end

  # Change the right subtree
  def set_right(node)
    if self.right
      self.right.parent = nil
    end
    self.right = node
    self.right.parent = self
    self
  end

  # Remove the child subtree from the matching node onwards
  def remove_child(node)
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
      return node
    end
    # If value is greater than current nodes value go right
    #   Recursively try to insert the same value on the right subtree if a right subtree exists
    #   Otherwise create a new node and add it as the right child of current node
    if val > self.value
      return self.right.insert(val) if self.right
      node = BinarySearchTreeNode.new
      self.set_right(node)
      return node
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
    node = self.find(val)
    return false unless node
    # If no left or rigth subtree
    #   If parent is there, remove this child from the parent
    #   Else set the value of this node to nil
    if !node.left && !node.right
      if node.parent
        node.parent.remove_child(node)
      else
        node.set_value(nil)
      end
    # If there is a left and a right subtree
    #   If the next highest value is not the right child, ie the right child has a left subtree
    #     Then find the min value node from the right node onwards
    #     And then make the node to be deleted the node with the min value from above search
    #       ie, recursively delete the next big node
    #       and make the value of the current node as the next big node
    #   If the next highest value node is the right child
    #     Then make this node the right node
    #     And make the right subtree of this node the right subtree of the right child
    elsif node.left && node.right
      next_big_node = node.right.find_min
      if next_big_node.value != node.right.value
        self.delete(next_big_node.value)
        node.set_value(next_big_node.value)
      else
        node.set_value(node.right.value)
        node.set_right(node.right.right)
      end
    # If there is only one subtree, ie, left or right
    #   Replace the current node with the node from the left or right subtree, which ever it is
    else
      if node.left
        node.parent.replace_child(node, node.left)
      else
        node.parent.replace_child(node, node.right)
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
