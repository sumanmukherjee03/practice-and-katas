# @param {String[]} words
# @param {Integer} k
# @return {String[]}
def top_k_frequent(words, k)
    h = {}
    words.each do |w|
      if h.has_key?(w)
        h[w] += 1
      else
        h[w] = 1
      end
    end

    h1 = {}
    h.each do |w,f|
      if h1.has_key?(f)
        h1[f] << w
      else
        h1[f] = [w]
      end
    end

    res = h1.keys.sort {|a,b| b <=> a}.map {|x| h1[x].sort}.flatten[0...k]

    return res
end
