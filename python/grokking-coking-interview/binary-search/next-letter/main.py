def describe():
    desc = """
Problem : Given an array of lowercase letters sorted in ascending order, find the smallest letter in the given array greater than a given key.
Assume the given array is a circular list, which means that the last letter is assumed to be connected with the first letter.
This also means that the smallest letter in the given array is greater than the last letter of the array and is also the first letter of the array.

Example :
    Input: ['a', 'c', 'f', 'h'], key = 'f'
    Output: 'h'
    Explanation: The smallest letter greater than 'f' is 'h' in the given array.


    Input: ['a', 'c', 'f', 'h'], key = 'b'
    Output: 'c'
    Explanation: The smallest letter greater than 'b' is 'c'.

    Input: ['a', 'c', 'f', 'h'], key = 'm'
    Output: 'a'
    Explanation: As the array is assumed to be circular, the smallest letter greater than 'm' is 'a'.
    """
    print(desc)

#  This follows a similar approach to the binary search approach for finding the ceiling of a number.
#  However, the array is circular. So consider the case where key is greater than or equal to the last letter or smaller than the first letter.
#  The other case is that if the exact key is found, then we want to find the next letter and not exactly that letter.
def search_next_letter(letters, key):
    if key >= letters[len(letters) - 1] or key < letters[0]:
        return letters[0]
    start, end = 0, len(letters) - 1
    while start <= end:
        mid = start + (end-start)//2
        if letters[mid] == key:
            return letters[mid+1]
        elif key < letters[mid]:
            end = mid-1
        elif key > letters[mid]:
            start = mid+1
    # since the loop is running until start <= end, so at the end of the while loop, start == end+1
    # we are not able to find the element in the given array, so the next big number will be arr[start]
    return letters[start]


def improved_search_next_letter(letters, key):
    n = len(letters)
    start, end = 0, n - 1
    while start <= end:
        mid = start + (end - start) // 2
        if key < letters[mid]:
            end = mid - 1
        else: # key >= letters[mid]:
            start = mid + 1
    #  since the loop is running until 'start <= end', so at the end of the while loop, 'start == end+1'
    #  in normal cases the modulo operator doesnt do anything. but in case that the search stops at the last letter,
    #  then start will be equal to len(letters), ie n - 1 + 1 == n. In that case, the 0th element will be the next big letter.
    #  That's when the modulo comes into play
    return letters[start % n]

#  Time complexity is O(log n)
def main():
    describe()

    input = ['a', 'c', 'f', 'h']
    key = 'f'
    print("Input : " + str(input) + " , " + key)
    print(search_next_letter(input, key))

    input = ['a', 'c', 'f', 'h']
    key = 'b'
    print("Input : " + str(input) + " , " + key)
    print(search_next_letter(input, key))

    input = ['a', 'c', 'f', 'h']
    key = 'm'
    print("Input : " + str(input) + " , " + key)
    print(search_next_letter(input, key))

main()
