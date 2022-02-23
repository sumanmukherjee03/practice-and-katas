def describe():
    desc = """
Problem : Given an array find the length of the longest consecutive sequence. Note that the consecutive sequence does not means that the numbers are present after each other in the array.

Example :
    Input : [14,1,8,4,0,13,6,9,2,7]
    Output : 4 because longest consecutive sequence is 6 7 8 9

    Input : [4,5,2,6,5,4,1,-5,0,4]
    Output : 3 because the longest consecutive sequence is 4 5 6

--------------------
    """
    print(desc)

#  Time complexity is O(n^3)
def longest_consecutive_sequence(arr):
    if len(arr) < 2:
        return len(arr)
    maxlen = 1
    for elm in arr:
        left = elm
        while left-1 in arr:
            left -= 1
        right = elm
        while right+1 in arr:
            right += 1
        maxlen = max(maxlen, right-left+1)
    return maxlen

#  A better solution is to first sort the array and then try and find the consecutive elements
#  Time complexity is O(n(log n)) but remember we are modifying the input
def longest_consecutive_sequence_by_sorting(arr):
    if len(arr) < 2:
        return len(arr)
    arr.sort()
    maxlen, current_maxlen = 1, 1
    for i in range(1, len(arr)):
        if arr[i] == arr[i-1]+1:
            current_maxlen += 1
        if arr[i] == arr[i-1]:
            pass
        else:
            current_maxlen = 1
        maxlen = max(maxlen, current_maxlen)
    return maxlen


#  Improving upon the first solution even further
#  Time complexity is O(n^2)
def longest_consecutive_sequence_slightly_optimized(arr):
    if len(arr) < 2:
        return len(arr)
    #  We put the elements of the array into a set, meaning no dupes.
    #  This reduces the time complexity when searching for contiguous elements because we dont have to traverse the array left and right
    #  Lookup in a set is O(1)
    values = set(arr)
    maxlen = 1
    for elm in arr:
        left = elm
        while left-1 in values:
            left -= 1
        right = elm
        while right+1 in values:
            right += 1
        maxlen = max(maxlen, right-left+1)
    return maxlen






#  We can improve the previous solution even more. Right now, we end up finding the same consecutive subsequence
#  for each element that is there in the subsequence. ie, in the first example, we end up finding the same subsequence once for each of 6, 7, 8 and 9.
#  This can be improved by only searching for the subsequence if there is not a smaller value in the set of values.
#  That would mean the current number would be the samllest number in the consecutive subsequence.
#  Although there is a nested loop, the nested loop would max run for n iterations only.
#  This gives us a time complexity of 2n, ie O(n)
def longest_consecutive_sequence_optimized(arr):
    if len(arr) < 2:
        return len(arr)
    #  We put the elements of the array into a set, meaning no dupes.
    #  This reduces the time complexity when searching for contiguous elements because we dont have to traverse the array left and right
    #  Lookup in a set is O(1)
    values = set(arr)
    maxlen = 1
    for elm in arr:
        if elm-1 in values:
            continue
        else:
            right = elm
            while right+1 in values:
                right += 1
            maxlen = max(maxlen, right-elm+1)
    return maxlen





def main():
    describe()

    input = [14,1,8,4,0,13,6,9,2,7]
    print("Input : " + str(input))
    print("Output : " + str(longest_consecutive_sequence_optimized(input)))

    input = [4,5,2,6,5,4,1,-5,0,4]
    print("Input : " + str(input))
    print("Output : " + str(longest_consecutive_sequence_optimized(input)))

main()
