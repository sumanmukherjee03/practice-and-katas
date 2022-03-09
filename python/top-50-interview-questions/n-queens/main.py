def describe():
    desc = """
Problem : Find the number of ways possible to place n queens on a n x n chess board such that no 2 queens can attack each other.
          The queens attack each other when they are placed on the same row, same column or placed diagonal to each other.

Example :
    Input : n = 4
    Output : 2


--------------------

    """
    print(desc)






#  The approach to this problem is to choose a place where a queen can be placed in each row.
#  So, we can keep going 1 queen per row and find possible spots.
#      - For possible spots start going column by column in each row.
#      - Place a queen in each column and check if it is a feasible solution. If yes, continue to next row. Else backtrack.
#  As we move on to next rows, it is possible that due to the placement of queens in the previous rows, all possible columns in the current row will get attacked.
#  In this case we backtrack to the previous row and try the next column in the previous row.
#
#  Time complexity is O(n^2 * n!)
def nQueens(n):
    board = [['.'] * n for i in range(n)]
    return _nQueens(n, board, 0)

def _nQueens(n, board, row):
    #  If we have reached past the last row then we have already found a solution, so return back 1 -> because you have found 1 way to put n queens on the board
    if row >= n:
        return 1
    #  Otherwise iterate over all the cells in the row to find possible solutions
    #  Once you find a solution add it to sumWays and then remove the queen from that row, because you want to continue exploring other solutions
    #  ie this is the part where you get to backtrack
    sumWays = 0
    for col in range(n):
        #  If placing the queen in the current cell is safe, then explore this solution further by considering next rows.
        #  Regardless of whether a solution is found or not remove the queen from that cell and explore the next cell in the row.
        if isNotAttacked(board, row, col):
            board[row][col] = 'Q'
            sumWays += _nQueens(n, board, row+1)
            board[row][col] = '.' # Remove the queen after finding possible solutions because yiou want to exp[lore solutions by placing queen in another cell
    return sumWays

#  How to check if a queen is attacked.
#      - We are not putting 2 queens in the same row, so we dont have to check the rows.
#      - We havent placed anything in the board beyond the current row, so, we dont need to check the bottom half of the board.
#      - We only need checking the column and the 2 diagonals and that too only on the part of the board above the current row.
#  The time complexity of this section is O(n)
def isNotAttacked(board, row, col):
    i = row-1 # Represents the row above the current row we want to check
    jLeft = col-1 # On the ith row, jLeft represents the cell diagonal from the current cell in consideration
    jRight = col+1 # On the ith row, jRight represents the cell in the other diagonal from the current cell in consideration

    #  Starting with the previous row of the current row, keep going up a row and keep checking the column in the rows above and the diagonal cells in the rows above
    while i >= 0:
        if board[i][col] == 'Q' or (jLeft >= 0 and board[i][jLeft] == 'Q') or (jRight < len(board) and board[i][jRight] == 'Q'):
            return False
        else:
            i -= 1
            jLeft -= 1
            jRight += 1
    return True







def main():
    describe()

    print("Input : ", 4)
    print("Output : ", nQueens(4))
    print("\n\n")

    print("Input : ", 8)
    print("Output : ", nQueens(8))
    print("\n\n")

main()
