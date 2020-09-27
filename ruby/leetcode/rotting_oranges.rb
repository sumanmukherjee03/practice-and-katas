# @param {Integer[][]} grid
# @return {Integer}
def oranges_rotting(grid)
  return -1 if grid.length == 0
  return 0 if grid.length == 1 && grid[0].length == 1 && grid[0][0] == 0
  return -1 if grid.length == 1 && grid[0].length == 1 && grid[0][0] == 1
  return 0 if grid.length == 1 && grid[0].length == 1 && grid[0][0] == 2

  mins = 0
  rows = grid.length
  cols = grid[0].length

  total_rotten_count = 0
  total_orange_count = 0
  initial_rotten = []

  i = 0
  while i < rows do
    j = 0
    while j < cols do
      total_orange_count += 1 if grid[i][j] > 0
      if grid[i][j] == 2
        initial_rotten << [i,j]
      end
      j += 1
    end
    i += 1
  end

  total_rotten_count += initial_rotten.length

  bfs = lambda do |rotten_collection|
    return nil if rotten_collection.length == 0
    # puts "--------mins : #{mins}, rotten collection : #{rotten_collection}----------"
    m = 0
    new_rotten = []
    while m < rotten_collection.length do
      i = rotten_collection[m][0] # row num
      j = rotten_collection[m][1] # col num
      # puts "processing neighbours of (#{i},#{j})"
      if i-1 >= 0 && grid[i-1][j] == 1
        grid[i-1][j] = 2
        new_rotten << [i-1,j]
      end
      if i+1 < rows && grid[i+1][j] == 1
        grid[i+1][j] = 2
        new_rotten << [i+1,j]
      end
      if j-1 >= 0 && grid[i][j-1] == 1
        grid[i][j-1] = 2
        new_rotten << [i,j-1]
      end
      if j+1 < cols && grid[i][j+1] == 1
        grid[i][j+1] = 2
        new_rotten << [i,j+1]
      end
      m += 1
    end

    mins += 1
    total_rotten_count += new_rotten.length
    # puts "mins : #{mins}, new rotten : #{new_rotten}"
    bfs.call(new_rotten)
  end

  bfs.call(initial_rotten)

  return total_rotten_count == total_orange_count ? mins-1 : -1
end
