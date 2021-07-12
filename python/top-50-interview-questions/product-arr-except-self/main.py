#  Problem : Given an array of integers write a function to return an array which
#  at arr[i] contains the product of all elements except element at index i
#  For ex : With input array [2,5,3,4]
#  Output : [60,24,40,30]


#  Precompute the products from the left and from the right, ie cumulative products from left and right.
#  Time complexity is O(n)
def productExceptSelf(arr):
    if len(arr) == 0:
        return []
    if len(arr) == 1:
        return [1]

    leftCumulativeProduct = [1] * len(arr)
    leftCumulativeProduct[0] = arr[0]
    for i in range(1, len(arr)):
        leftCumulativeProduct[i] = leftCumulativeProduct[i-1] * arr[i]

    rightCumulativeProduct = [1] * len(arr)
    rightCumulativeProduct[len(arr)-1] = arr[len(arr)-1]
    for i in range(len(arr)-2,-1,-1):
        rightCumulativeProduct[i] = rightCumulativeProduct[i+1] * arr[i]

    out = [1] * len(arr)
    for i in range(len(arr)):
        p = 1
        if i == 0:
            p *= rightCumulativeProduct[i+1]
        elif i == len(arr)-1:
            p *= leftCumulativeProduct[i-1]
        else:
            p *= (leftCumulativeProduct[i-1] * rightCumulativeProduct[i+1])
        out[i] = p

    return out
