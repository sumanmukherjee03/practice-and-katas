# @param {Integer[][]} intervals
# @return {Integer[][]}
def merge(intervals)
  # [[1,3],[2,6],[8,10],[15,18]]
  sets = []
  i = 0
  while i < intervals.length do
    l = intervals[i].first
    r = intervals[i].last

    j = 0
    set_union_happened = false
    while j < sets.length do
      l_included = false
      r_included = false
      set_l_included = false
      set_r_included = false

      if l >= sets[j].first && l <= sets[j].last
        l_included = true
      end
      if r >= sets[j].first && r <= sets[j].last
        r_included = true
      end

      if sets[j].first >= l && sets[j].first <= r
        set_l_included = true
      end
      if sets[j].last >= l && sets[j].last <= r
        set_r_included = true
      end

      if l_included || r_included
        if l_included && !r_included
          sets[j][1] = r
          set_union_happened = true
        elsif !l_included && r_included
          sets[j][0] = l
          set_union_happened = true
        else
          set_union_happened = true
        end
      end

      if set_l_included || set_r_included
        if set_l_included && !set_r_included
          sets[j][0] = l
          set_union_happened = true
        elsif !set_l_included && set_r_included
          sets[j][1] = r
          set_union_happened = true
        else
          sets[j][0] = l
          sets[j][1] = r
          set_union_happened = true
        end
      end

      break if set_union_happened
      j += 1
    end

    sets << intervals[i] unless set_union_happened
    i += 1
  end

  if intervals.length > sets.length
    return merge(sets)
  else
    return sets
  end
end
