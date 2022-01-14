def describe():
    desc = """
Problem : Given an array, find the average of all subarrays of K contiguous elements in it.
     Ex : Given Array: [1, 3, 2, 6, -1, 4, 1, 8, 2], K=5
     Output: [2.2, 2.8, 2.4, 3.6, 2.8]
    """
    print(desc)

#  Time complexity = O(n^2)
def brute_force_find_averages_of_subarrays(K, arr):
    result = []
    print("Length of arr is : {}, K is : {}, Var i will iterate upto (not including) : {}".format(len(arr), K, len(arr)-K+1))
    for i in range(len(arr)-K+1):
        # find sum of next 'K' elements
        _sum = 0.0
        for j in range(i, i+K):
            _sum += arr[j]
        result.append(_sum/K)  # calculate average
    return result


#  Time complexity = O(n)
#  Maintain a container for collecting the result
#  Maintain a pointer for the start of the window and one for the end of the window
#  Move end pointer from start one by one UNTIL it is of the length of the sliding window
#  Calculate aggregate, and insert into the results collection
#  Move end pointer by one element right and start pointer by one element right
def find_averages_of_subarrays(K, arr):
    result = []
    windowSum, windowStart = 0.0, 0
    for windowEnd in range(len(arr)):
        # add the next element - keep adding the next element until you hit the length of the sliding window
        windowSum += arr[windowEnd]
        # slide the window, we don't need to slide if we've not hit the required window size of 'k'
        #  ie, the index of the end pointer of the window should be start of window + K or more.
        #  if windowEnd >= windowStart + K - 1:
        if windowEnd >= windowStart + K - 1:
            result.append(windowSum / K)  # calculate the average
            windowSum -= arr[windowStart]  # subtract the element going out
            windowStart += 1  # slide the window ahead
    return result


def main():
    describe()
    result1 = brute_force_find_averages_of_subarrays(5, [1, 3, 2, 6, -1, 4, 1, 8, 2])
    print("Averages of subarrays of size K (solution 1): " + str(result1))
    result2 = find_averages_of_subarrays(5, [1, 3, 2, 6, -1, 4, 1, 8, 2])
    print("Averages of subarrays of size K (solution 2): " + str(result2))


main()
