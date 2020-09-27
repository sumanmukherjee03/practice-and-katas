# Definition for a binary tree node.
# class TreeNode
#     attr_accessor :val, :left, :right
#     def initialize(val)
#         @val = val
#         @left, @right = nil, nil
#     end
# end

# Encodes a tree to a single string.
#
# @param {TreeNode} root
# @return {string}
def serialize(root)
  res = []
  return res unless root

  # Maintain a queue for BFS to take care of the nodes at each level
  q = Queue.new
  q.push(root)

  while !q.empty? do
    i = 0
    nodes_in_level = []
    no_nodes_in_queue = q.size
    while i < no_nodes_in_queue do
      node = q.pop()

      if node
        nodes_in_level << node.val
      else
        nodes_in_level << "null"
      end

      if node && node.left
        q.push(node.left)
      elsif node && !node.left
        q.push(nil)
      end

      if node && node.right
        q.push(node.right)
      elsif node && !node.right
        q.push(nil)
      end
      i += 1
    end
    res = res.concat(nodes_in_level)
  end

  # Remove the trailing nils
  idx = res.rindex { |m| m != "null" }
  res = res[0..idx]
  puts "Final res : #{res}"
  return res
end

# Decodes your encoded data to tree.
#
# @param {string} data
# @return {TreeNode}
def deserialize(data)
  # For deserializing a level order tree traversal
  level_order_traversal = lambda do |i|
    idx = i + 1 # Imagine as if the indexing is starting from 1
    node = TreeNode.new(data[i])
    if (idx * 2 - 1) < data.length
      node.left = level_order_traversal.call(idx * 2 - 1)
    end
    if (idx * 2 - 1 + 1) < data.length
      node.right = level_order_traversal.call(idx * 2 - 1 + 1)
    end

    return node
  end

  root = level_order_traversal.call(0)
  return root
end

# [1,2,3,null,null,4,5]

# Your functions will be called as such:
# deserialize(serialize(data))
