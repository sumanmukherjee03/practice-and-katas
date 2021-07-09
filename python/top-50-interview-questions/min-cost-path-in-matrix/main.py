#  Problem : Given a matrix of integers of size n * m, where each cell matrix[i][j]
#  represents the cost of passing through that cell, create a function that returns the
#  minimum cost of going from top left cell to bottom right cell.
#  Knowing that you can only move bottom or right.

#  Ex :
    #  3  12  4  7  10
    #  6  8  15 11  4
    #  19 5  14 32  21
    #  1  20  2  9  7

#  Output : 54 -> 3,6,8,5,14,2,9,7

#  Intuitively it might look like a greedy approach might work, but in reality it would not.


#  For each cell the minimum cost path is
#      the cost of the cell itself + whatever is the minimum of -> path going right or path going down
#      So, explore both options recursively
#  Time complexity of this solution is O(2^(rows+cols)) or O(2^(n+m))
def minimumCostPath01(matrix, i = 0, j = 0):
    if len(matrix) == 0:
        return 0

    rows = len(matrix)
    cols = len(matrix[0])
    # If your origin is the right bottom cell
    if i == rows-1 and j == cols-1:
        return matrix[i][j]
    # If your origin is on the last column
    #   return the cost of the current cell + minimum cost from cell below it to the destination
    #   because there is no other direction to traverse
    elif i < rows and j == cols-1:
        return matrix[i][j] + minimumCostPath01(matrix, i+1, j)
    # If your origin is on the last row
    #   return the cost of the current cell + minimum cost from cell right of it to the destination
    #   because there is no other direction to traverse
    elif i == rows-1 and j < cols:
        return matrix[i][j] + minimumCostPath01(matrix, i, j+1)
    #  If your origin is anywhere else
    else:
        return matrix[i][j] + min(minimumCostPath01(matrix, i+1, j), minimumCostPath01(matrix, i, j+1))



#  The solution above is quite expensive in terms of time complexity.
#  And mainly because if you trace it's recursion call stack there are many repeated operations.
#  This is a good fit for dynamic programming.
#  Time complexity is O(n*m)
def minimumCostPath02(matrix):
    if len(matrix) == 0:
        return 0

    rows = len(matrix)
    cols = len(matrix[0])

    # Create a 2D array to store the results of dynamic programming and initialize
    # the cost of starting from any cell i,j to the end as infinity
    dp = [[float("inf")]*cols for i in range(rows)]

    for i in range(rows-1,-1,-1):
        for j in range(cols-1,-1,-1):
            # If your origin is the right bottom cell
            if i == rows-1 and j == cols-1:
                dp[i][j] = matrix[i][j]
            # If your origin is on the last column
            #   return the cost of the current cell + minimum cost from cell below it to the destination
            #   because there is no other direction to traverse
            elif i < rows and j == cols-1:
                dp[i][j] = matrix[i][j] + dp[i+1][j]
            # If your origin is on the last row
            #   return the cost of the current cell + minimum cost from cell right of it to the destination
            #   because there is no other direction to traverse
            elif i == rows-1 and j < cols:
                dp[i][j] = matrix[i][j] + dp[i][j+1]
            # If your origin is anywhere else
            else:
                dp[i][j] = matrix[i][j] + min(dp[i+1][j], dp[i][j+1])

    return dp[0][0]




# Same solution for dynamic programming but instead of calculating from the back, lets calculate from the top
# Time complexity is O(n*m)
def minimumCostPath03(matrix):
    if len(matrix) == 0:
        return 0

    rows = len(matrix)
    cols = len(matrix[0])

    # Create a 2D array to store the results of dynamic programming and initialize
    # the cost of starting from any cell i,j to the end as infinity
    dp = [[0]*cols for i in range(rows)]

    dp[0][0] = matrix[0][0]

    #  fill the first row
    for j in range(1,cols):
        dp[0][j] = dp[0][j-1] + matrix[0][j]

    #  fill the first column
    for i in range(1,rows):
        dp[i][0] = dp[i-1][0] + matrix[i][0]

    # now using the prefilled row above and column on the left, fill the rest of the grid
    for i in range(1,rows):
        for j in range(1,cols):
            dp[i][j] = matrix[i][j] + min(dp[i-1][j],dp[i][j-1])

    return dp[rows-1][cols-1]
