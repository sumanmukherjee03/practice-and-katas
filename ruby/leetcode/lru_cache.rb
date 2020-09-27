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
            index = self.ordered_keys.index(key)
            self.ordered_keys = self.ordered_keys[0...index] + self.ordered_keys[index+1..-1] + [key]
            res = self.store[key]
        end
        res
    end


=begin
    :type key: Integer
    :type value: Integer
    :rtype: Void
=end
    def put(key, value)
        is_key_present = self.store.has_key?(key)
        invalidate unless is_key_present
        self.store[key] = value
        if !is_key_present
            self.ordered_keys << key
        else
            index = self.ordered_keys.index(key)
            self.ordered_keys = self.ordered_keys[0...index] + self.ordered_keys[index+1..-1] + [key]
            res = self.store[key]
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

end

# Your LRUCache object will be instantiated and called as such:
# obj = LRUCache.new(capacity)
# param_1 = obj.get(key)
# obj.put(key, value)
