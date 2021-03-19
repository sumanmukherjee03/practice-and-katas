# Design a data structure that follows the constraints of a Least Recently Used (LRU) cache.
# Implement the LRUCache class:

# LRUCache(int capacity) Initialize the LRU cache with positive size capacity.
# int get(int key) Return the value of the key if the key exists, otherwise return -1.
# void put(int key, int value) Update the value of the key if the key exists. Otherwise, add the key-value pair to the cache. If the number of keys exceeds the capacity from this operation, evict the least recently used key.
# Follow up:
# Could you do get and put in O(1) time complexity?

# Example 1:

# Input
# ["LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"]
# [[2], [1, 1], [2, 2], [1], [3, 3], [2], [4, 4], [1], [3], [4]]
# Output
# [null, null, null, 1, null, -1, null, -1, 3, 4]

# Explanation
# LRUCache lRUCache = new LRUCache(2);
# lRUCache.put(1, 1); // cache is {1=1}
# lRUCache.put(2, 2); // cache is {1=1, 2=2}
# lRUCache.get(1);    // return 1
# lRUCache.put(3, 3); // LRU key was 2, evicts key 2, cache is {1=1, 3=3}
# lRUCache.get(2);    // returns -1 (not found)
# lRUCache.put(4, 4); // LRU key was 1, evicts key 1, cache is {4=4, 3=3}
# lRUCache.get(1);    // return -1 (not found)
# lRUCache.get(3);    // return 3
# lRUCache.get(4);    // return 4


# Constraints:

# 1 <= capacity <= 3000
# 0 <= key <= 3000
# 0 <= value <= 104
# At most 3 * 104 calls will be made to get and put
#

class LRUCache
  attr_accessor :capacity, :store, :ordered_keys

=begin
    :type capacity: Integer
=end
  def initialize(capacity)
    @capacity = capacity
    @store = {}
    @ordered_keys = []
  end


=begin
    :type key: Integer
    :rtype: Integer
=end
  def get(key)
    is_key_present = self.store.has_key?(key)
    res = -1
    if is_key_present
      res = move_elm_to_end(key)
    end
    res
  end


=begin
    :type key: Integer
    :type value: Integer
    :rtype: Void
=end
  # If key is not present check the capacity vs length of the ordered list and if required remove the oldest element in the queue
  # The oldest element is always the first element in the queue
  # If key is not present then push the new key at the end of the ordered list which here is acting as a queue
  # Otherwise remove the key from the position it was at in the ordered list and push it to the end of the queue
  def put(key, value)
    is_key_present = self.store.has_key?(key)
    invalidate unless is_key_present
    self.store[key] = value
    if !is_key_present
      self.ordered_keys << key
    else
      move_elm_to_end(key)
    end
    self
  end

  def invalidate
    if self.ordered_keys.length == capacity
      k = self.ordered_keys.first
      self.store.delete(k)
      self.ordered_keys = self.ordered_keys[1..-1] # delete first element
    end
    self
  end

  def move_elm_to_end(key)
    index = self.ordered_keys.index(key)
    self.ordered_keys = self.ordered_keys[0...index] + self.ordered_keys[index+1..-1] + [key]
    return self.store[key]
  end
end

# Your LRUCache object will be instantiated and called as such:
# obj = LRUCache.new(capacity)
# param_1 = obj.get(key)
# obj.put(key, value)
