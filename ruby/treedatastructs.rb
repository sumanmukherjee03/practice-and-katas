class BinaryTreeNode
  attr_accessor :left, :right, :parent, :value

  def initialize(value = nil)
    @left = nil
    @right = nil
    @parent = nil
    @value = value
  end

  def left_height
    return 0 unless self.left
    self.left.height + 1
  end

  def right_height
    return 0 unless self.right
    self.right.height + 1
  end

  def height
    [self.left_height, self.right_height].max
  end

  def balance
    self.left_height - self.right_height
  end

  def set_value(val)
    self.value = val
    self
  end

  def set_left(node)
    if self.left
      self.left.parent = nil
    end
    self.left = node
    self.left.parent = self
    self
  end

  def set_right(node)
    if self.right
      self.right.parent = nil
    end
    self.right = node
    self.right.parent = self
    self
  end

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

  def inorder
    traversal = []
    traversal.concat(self.left.inorder) if self.left
    traversal << self.value
    traversal.concat(self.right.inorder) if self.right
    traversal
  end

  def preorder
    traversal = []
    traversal << self.value
    traversal.concat(self.left.preorder) if self.left
    traversal.concat(self.right.preorder) if self.right
    traversal
  end

  def postorder
    traversal = []
    traversal.concat(self.left.postorder) if self.left
    traversal.concat(self.right.postorder) if self.right
    traversal << self.value
    traversal
  end

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

  def self.copy(src, dest)
    dest.set_value(src.value)
    dest.set_left(src.left)
    dest.set_right(src.right)
  end
end

class BinarySearchTreeNode < BinaryTreeNode
  def insert(val)
    return self unless val
    if val < self.value
      return self.left.insert(val) if self.left
      node = BinarySearchTreeNode.new
      self.set_left(node)
      return node
    end
    if val > self.value
      return self.right.insert(val) if self.right
      node = BinarySearchTreeNode.new
      self.set_right(node)
      return node
    end
    self
  end

  def find(val)
    return self if self.value == val
    return self.left.find(val) if val < self.value && self.left
    return self.right.find(val) if val > self.value && self.right
    return nil
  end

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
    #     And then make the node to be deleted the node with the min value from above
    #     ie, recursively delete the min value node from above
    #     and make the value of the current node as the node from above
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

  def find_min
    return self unless self.left
    self.left.find_min
  end
end

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
