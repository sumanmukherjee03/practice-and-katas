#  Problem : Given a matrix of size n*m where each cell contains a value of 0 or 1
#  where 0 means cell is empty and 1 means cell has a wall, create a function that
#  returns the numbers of paths that we can take to go from top left cell to bottom right cell
#  knowing that you can only move to the right or the bottom.


#  This solution is to recursively explore the 2 paths going left and right and if there is no wall then keep exploring
#  If the bottom right cell is reached as part of the exploration you increase the count of paths possible
def paths01(matrix):
    if len(matrix) == 0:
        return 0
    if len(matrix) == 1 and len(matrix[0]) == 1:
        if matrix[0][0] == 0:
            return 1
        else:
            return 0

    #  Maintain the count in a global variable that we reference as a list inside the inner function
    numPaths = [0]

    def explore(row, col):
        if matrix[row][col] == 1:
            return

        bottomRightCornerRowIndex = len(matrix) - 1
        bottomRightCornerColIndex = len(matrix[0]) - 1

        if row == bottomRightCornerRowIndex and col == bottomRightCornerColIndex:
            numPaths[0] += 1
            return
        elif row == bottomRightCornerRowIndex and col < bottomRightCornerColIndex:
            explore(row, col+1)
        elif row < bottomRightCornerRowIndex and col == bottomRightCornerColIndex:
            explore(row+1, col)
        else:
            explore(row, col+1)
            explore(row+1, col)
        return

    explore(0,0,[])
    return numPaths[0]



#  This one is a recursive programming solution as well
#  Time complexity is O(2^(n+m))
def paths02(matrix, i = 0, j = 0):
    if len(matrix) == 0:
        return 0

    rows = len(matrix)
    cols = len(matrix[0])
    if i > rows - 1 or j > cols - 1 or matrix[i][j] == 1:
        return 0
    elif i == rows - 1 and j == cols - 1:
        return 1
    else:
        return paths02(matrix, i+1, j) + paths02(matrix, i, j+1)



#  This one is a dynamic programming solution
#  Time complexity is O(n*m)
def paths03(matrix):
    if len(matrix) == 0:
        return 0

    rows = len(matrix)
    cols = len(matrix[0])
    # Array dp stores the number of paths to destination (ie bottom right cell) when origin is cell (i,j)
    dp = [[0] * cols for i in range(rows)]

    for i in range(rows-1,-1,-1):
        for j in range(cols-1,-1,-1):
            if i == rows-1 and j == cols-1:
                if matrix[i][j] == 0:
                    dp[i][j] = 1
                else:
                    dp[i][j] = 0
            elif i == rows-1 and j < cols-1:
                if matrix[i][j] == 0:
                    dp[i][j] = dp[i][j+1]
                else:
                    dp[i][j] = 0
            elif i < rows-1 and j == cols-1:
                if matrix[i][j] == 0:
                    dp[i][j] = dp[i+1][j]
                else:
                    dp[i][j] = 0
            else:
                if matrix[i][j] == 0:
                    dp[i][j] = dp[i][j+1] + dp[i+1][j]
                else:
                    dp[i][j] = 0

    return dp[0][0]



#  This is another dynamic programming solution with the dp array different from the one above
#  The time complexity of this solution is O(n*m)
def paths04(matrix):
    if len(matrix) == 0:
        return 0

    rows = len(matrix)
    cols = len(matrix[0])

    #  Here dp stores the number of paths from the top left cell to current cell (i, j)
    dp = [[0] * cols for i in range(rows)]
    if matrix[0][0] == 0:
        dp[0][0] = 1
    else:
        dp[0][0] = 0

    #  top left cell
    #  ^
    #  |
    #  ^
    #  |
    #  (i,j)
    for i in range(1,rows):
        if matrix[i][0] == 0:
            dp[i][0] = dp[i-1][0]
        else:
            dp[i][0] = 0

    #  top left cell <- <- <- (i,j)
    for j in range(1,cols):
        if matrix[0][j] == 0:
            dp[0][j] = dp[0][j-1]
        else:
            dp[0][j] = 0

    #      ^
    #      |
    #  <- i,j
    for i in range(1,rows):
        for j in range(1,cols):
            if matrix[i][j] == 0:
                dp[i][j] = dp[i-1][j] + dp[i][j-1]
            else:
                dp[0][j] = 0

    #  This element in dp represents the number of paths from top right cell to the bottom right cell
    return dp[rows-1][cols-1]
