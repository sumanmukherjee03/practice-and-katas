# @param {Integer[][]} matrix
# @return {Integer[]}
def spiral_order(matrix)
  return matrix if matrix.empty?

  rows = matrix.length
  cols = matrix[0].length
  visited = []

  first_col = 0
  first_row = 0
  last_col = cols - 1
  last_row = rows - 1

  visit = lambda do |i,j|
    while visited.length < (rows * cols) do
      visited << matrix[i][j]

      case
      when i == first_row && j < last_col
        visit.call(i,j+1)
      when i == last_row && j > first_col
        visit.call(i,j-1)
      when j == last_col && i < last_row
        visit.call(i+1,j)
      when j == first_col && i > first_row + 1
        visit.call(i-1,j)
      else
        first_col += 1
        first_row += 1
        last_col -= 1
        last_row -= 1
        visit.call(i,j+1)
      end
    end
  end

  visit.call(0,0)

  return visited
end
