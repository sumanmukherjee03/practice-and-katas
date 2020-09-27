# @param {Integer[]} nums1
# @param {Integer[]} nums2
# @return {Float}
def find_median_sorted_arrays(nums1, nums2)
  l1 = nums1.length
  l2 = nums2.length
  if l1 > l2
    outer = nums1
    inner = nums2
  else
    outer = nums2
    inner = nums1
  end

  i=0
  j=0
  k = 0
  arr = []
  while i < outer.length do
    if j < inner.length && outer[i] < inner[j]
      arr[k] = outer[i]
      k += 1
    elsif j < inner.length && outer[i] >= inner[j]
      while j < inner.length && outer[i] >= inner[j] do
        arr[k] = inner[j]
        k += 1
        j += 1
      end
      arr[k] = outer[i]
      k += 1
    else
      arr[k] = outer[i]
      k += 1
    end
    i += 1
  end

  if j < inner.length
    arr = arr.concat(inner[j..-1])
  end

  mid = arr.length/2
  res = 0
  if arr.length % 2 == 1
    res = Float(arr[mid])
  else
    res = Float(arr[mid-1] + arr[mid])/2
  end

  return res
end
