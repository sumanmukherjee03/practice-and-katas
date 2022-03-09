from collections import deque

def describe():
    desc = """
Problem : Given a board of characters and a string word, check if you can find the word in the board.
          Words must be with adjacent cells, ie either horizontal or vertical.
          The same cell can be used only once.

Example :
    Input : [['K','I','N','T'], ['B','I','N','S'], ['G','N','Y','I'], ['U','O','E','D'], ['D','I','B','V'], ['H','I','R','T']], word = 'INSIDE'
    Output : true
    Explanation : Each of these inner arrays are rows of chars in the board. The pattern goes like this - board[1][1:4] + board[2][3] + board[3][3] + board[3][2]
                  Draw it out to get a better understanding of how it looks like.


    Input : [['K','I','N','T'], ['B','I','N','S'], ['G','N','Y','I'], ['U','O','E','D'], ['D','I','B','V'], ['H','I','R','T']], word = 'CODE'
    Output : false
    Explanation : The word 'CODE' does not exist in the board

----------------------
    """
    print(desc)



def wordExists(input, word):
    chars = list(word)
    noRows = len(input)
    noCols = len(input[0])

    #  First find all the starting points.
    #  May be one of these starting points will yield the entire word, but most of them wont.
    possibleStartingPoints = deque()
    for i in range(0,noRows):
        for j in range(0,noCols):
            if input[i][j] == chars[0]:
                possibleStartingPoints.append([i,j])

    #  Explore each starting point one at a time
    while possibleStartingPoints:
        r,c = possibleStartingPoints.popleft()
        q = deque() # Declare a queue so that we can explore the neighbours of this starting position
        q.append([r,c]) # Insert the first/starting position into the queue
        charPtr = 0 # Use a pointer to keep track of the next char to match for the word you want to find
        #  Keep iterating until either the entire word is found or the queue is empty and no more neighbours need to be explored
        #  This section is similar to DFS, just a slightly improvised version of it based on which neighbours we want to explore
        while q and charPtr < len(chars):
            charFound = False
            temp = 0
            noNeighboursToExplore = len(q) # At any point the number of items of the queue are tge neighbours to explore
            for _ in range(noNeighboursToExplore):
                i, j = q.popleft()
                if input[i][j] == chars[charPtr]:
                    print("<" + str(i) + "," + str(j) + ">")
                    charFound = True
                    if i+1 < noRows:
                        q.append([i+1,j])
                        temp += 1
                    if j+1 < noCols:
                        q.append([i,j+1])
                        temp += 1
                    if i-1 >= 0:
                        q.append([i-1,j])
                        temp += 1
                    if j-1 >= 0:
                        q.append([i,j-1])
                        temp += 1
            if charFound:
                charPtr += 1 # If the char we were looking for is found we move forward to search for the next char

        if charPtr == len(chars):
            return True

    return False





#  Time complexity of the recurring search function
#       T(n, m, 0) = 1
#       T(n, m, w) = 4T(n, m, w-1) + 1
#       .... and so on
#       Eventually following this recurrence formula and induction we will get a time complexity of O(4^w) for the search function
#  So, total time complexity of the solution is O(nm(4^w))
def wordExistsConcise(board, word):
    noRows = len(board)
    noCols = len(board[0])
    visited = set()
    for i in range(noRows):
        for j in range(noCols):
            if board[i][j] == word[0]:
                if search(board, word, i, j, 0, visited):
                    return True
    return False

def outOfBoard(board, i, j):
    noRows = len(board)
    noCols = len(board[0])
    return i >= noRows or j >= noCols

def search(board, word, i, j, counter, visited):
    #  If you have reached the end of the word then you have found it, so return true
    if counter == len(word):
        return True
    #  Return false in 3 cases
    #  - the cell we are referring to is out of the board
    #  - the value of the cell on the board is not the char in the word we are matching against
    #  - this cell has already been visited
    #     - This is to take care of the condition that the same cell can not be used twice when matching the word
    #       Otherwise it is quite possible that you might come back to that same cell again as part of the traversal and get caught up in an infinite loop
    elif outOfBoard(board, i, j) or (i, j) in visited or board[i][j] != word[counter]:
        return False
    else:
        #  If we have reached this stage that means we have finally had a match with the character of the word we want to match and this cell in the board
        #  At this point we mark the cell as visited and explore it's 4 neighbours recursively
        #  If after recursive exploration of the subproblem with the 4 neighbours, one of them end up returning true, that means we found the word
        #  Otherwise we did not fnd it, so remove this cell from visited and return false. We remove this cell from visited set because we only found a partial match with the word.
        #  This is the same as tracking back in the recursive call stack. Untimately, the visited set will become empty if only a partial match of the word was found.
        #  If a match was indeed found then visited contains the actual path/solution
        visited.add((i, j))
        if search(board, word, i+1, j, counter+1, visited) or search(board, word, i, j+1, counter+1, visited) or search(board, word, i-1, j, counter+1, visited) or search(board, word, i, j-1, counter+1, visited):
            return True
        else:
            visited.remove((i,j))
            return False







def main():
    describe()

    input = [['K','I','N','T'], ['B','I','N','S'], ['G','N','Y','I'], ['U','O','E','D'], ['D','I','B','V'], ['H','I','R','T']]
    word = 'INSIDE'
    print("Input : ")
    for r in input:
        print(r)
    print("\n")
    print(word)
    print("\n")
    print("Output : " + str(wordExistsConcise(input, word)))
    print("\n\n")

    input = [['K','I','N','T'], ['B','I','N','S'], ['G','N','Y','I'], ['U','O','E','D'], ['D','I','B','V'], ['H','I','R','T']]
    word = 'CODE'
    print("Input : ")
    for r in input:
        print(r)
    print("\n")
    print(word)
    print("\n")
    print("Output : " + str(wordExists(input, word)))
    print("\n\n")

main()
