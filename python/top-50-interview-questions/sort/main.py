#  Problem : Sort an array

# Time complexity O(n^2)
def bubbleSort(arr):
    if len(arr) <= 1:
        return arr
    #  Run an outer loop len(arr) times
    i = 0
    while i < len(arr):
        #  Keep swapping adjacent elements, ie sort adjacent elements in the inner loop
        j = 1
        while j < len(arr):
            if arr[j-1] > arr[j]:
                arr[j-1], arr[j] = arr[j], arr[j-1]
            j += 1
        i += 1
    return arr



# Time complexity of merge sort is O(nlogn)
def mergeSort(arr):
    if len(arr) <= 1:
        return arr

    def merge(left, right):
        if len(left) == 0:
            return right
        if len(right) == 0:
            return left
        res = []
        i = 0
        j = 0
        while i < len(left) and j < len(right):
            if left[i] <= right[j]:
                res.append(left[i])
                i += 1
            else:
                res.append(right[j])
                j += 1
        if i < len(left):
            res = res + left[i:len(left)]
        if j < len(right):
            res = res + right[j:len(right)]
        return res


    mid = len(arr)//2
    l = mergeSort(arr[0:mid])
    r = mergeSort(arr[mid:len(arr)])

    return merge(l, r)
