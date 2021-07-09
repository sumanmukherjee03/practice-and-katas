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



#  For each cell the minimum cost path is
#      the cost of the cell itself + whatever is the minimum of -> path going right or path going down
#      So, explore both options recursively
def minimumCostPath(matrix, i = 0, j = 0):
    if len(matrix) == 0:
        return 0

    rows = len(matrix)
    cols = len(matrix[0])
    # If you are on the last cell
    if i == rows-1 and j == cols-1:
        return matrix[i][j]
    # If you are on the last column
    #   return the cost of the current cell + minimum cost from cell below it to the destination
    #   because there is no other direction to traverse
    elif i < rows and j == cols-1:
        return matrix[i][j] + minimumCostPath(matrix, i+1, j)
    # If you are on the last row
    #   return the cost of the current cell + minimum cost from cell right of it to the destination
    #   because there is no other direction to traverse
    elif i == rows-1 and j < cols:
        return matrix[i][j] + minimumCostPath(matrix, i, j+1)
    else:
        return matrix[i][j] + min(minimumCostPath(matrix, i+1, j), minimumCostPath(matrix, i, j+1))
