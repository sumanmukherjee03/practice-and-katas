# Given an m x n 2D binary grid grid which represents a map of '1's (land) and '0's (water), return the number of islands.

# An island is surrounded by water and is formed by connecting adjacent lands horizontally or vertically. You may assume all four edges of the grid are all surrounded by water.

# Example 1:

# Input: grid = [
  # ["1","1","1","1","0"],
  # ["1","1","0","1","0"],
  # ["1","1","0","0","0"],
  # ["0","0","0","0","0"]
# ]
# Output: 1
# Example 2:

# Input: grid = [
  # ["1","1","0","0","0"],
  # ["1","1","0","0","0"],
  # ["0","0","1","0","0"],
  # ["0","0","0","1","1"]
# ]
# Output: 3


# Constraints:

# m == grid.length
# n == grid[i].length
# 1 <= m, n <= 300
# grid[i][j] is '0' or '1'

# @param {Character[][]} grid
# @return {Integer}
def num_islands(grid)
  count = 0
  rows = grid.length
  cols = grid[0].length if rows > 0

  # visit all the neighbours and keep converting the currently visited node to "0" from "1"
  bfs = lambda do |r, c|
    return if grid[r][c] == "0"
    grid[r][c] = "0"
    bfs.call(r+1, c) if r+1 < rows
    bfs.call(r, c+1) if c+1 < cols
    bfs.call(r-1, c) if r-1 >= 0
    bfs.call(r, c-1) if c-1 >= 0
  end

  # At least one "1" represents an island. As soon as you find a "1" start exploring it's neighbours
  # And as you visit the neighbours start converting them to "0" to mark as visited
  rows.times do |i|
    cols.times do |j|
      if grid[i][j] == "1"
        count += 1
        bfs.call(i, j)
      end
    end
  end

  return count
end

