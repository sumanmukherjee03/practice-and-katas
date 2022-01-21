def describe():
    desc = """
Problem : Given a sorted array, create a new array containing squares of all the numbers of the input array in the sorted order.
Example :
    Input: [-2, -1, 0, 2, 3]
    Output: [0, 1, 4, 4, 9]
    """
    print(desc)

#  This problem becomes a little hard because of the negative numbers and the fact that squares of negative numbers can be bigger than that of the square of positive nums.
#  Find the point in the given sorted array where numbers transition from negative to positive.
#  Then traverse with 2 pointers, one forward and one backward. Compare the squares and insert the smaller square into the resulting array.
#  Move the pointers forward or backward accordingly.
#  Time complexity is O(n)
def make_squares(arr):
    res = []
    non_neg_index = -1
    for i in range(0, len(arr)):
        if arr[i] < 0:
            non_neg_index = i
    if non_neg_index >= 0:
        j = non_neg_index+1
        k = non_neg_index
        large_val = 9999999 # For some reason math.inf did not work here possibly due to some type issue
        while j < len(arr) or k >= 0:
            n1 = arr[j]**2 if j < len(arr) else large_val
            n2 = arr[k]**2 if k >= 0 else large_val
            if n1 > n2:
                res.append(n2)
                k -= 1
            else:
                res.append(n1)
                j += 1
    else:
        for i in range(0, len(arr)):
            res.append(arr[i]**2)
    return res

#  In this solution we maintain 2 pointers, one goes left from 0 and the other traverses back from the end.
#  Compare squares and insert into resulting array.
def alternate_make_squares(arr):
  n = len(arr)
  squares = [0 for x in range(n)]
  highestSquareIdx = n - 1
  left, right = 0, n - 1
  while left <= right:
    leftSquare = arr[left] * arr[left]
    rightSquare = arr[right] * arr[right]
    if leftSquare > rightSquare:
      squares[highestSquareIdx] = leftSquare
      left += 1
    else:
      squares[highestSquareIdx] = rightSquare
      right -= 1
    highestSquareIdx -= 1
  return squares

def main():
    describe()
    arr = [-2, -1, 0, 2, 3]
    res = make_squares(arr)
    print("Input : ", arr)
    print("Output : ", res)

main()
