def describe():
    desc = """
Problem :
    Given two integer arrays to represent weights and profits of N items, we need to find a subset of these items
    which will give us maximum profit such that their cumulative weight is not more than a given number C.
    Each item can only be selected once, which means we either put that item in the knapsack or we skip it.
    """
    print(desc)

#  Recursive solution :
#       for each item i
#           explore a solution which includes the item i
#           recursively process the remaining capacity and the rest of the items
#       explore a solution which does not include the item i
#            recursively process the remaining capacity and the rest of the items
#       return the solution set with the highest profit
def solve_knapsack(profits, weights, capacity):
    #  Create a 2d array with each cell initialized to -1.
    #  The rows represent index of the element being considered in the solution set.
    #  The columns here represent each capacity broken down by an integer value.
    #  And a cell represents a recursive call for a <index,capacity> combination.
    #  As we will be solving the subproblems for a capacity value and index, we can store them in the dp array and later use those values
    dp = [[-1 for x in range(capacity+1)] for y in range(len(profits))]
    return solve_knapsack_recursive(dp, profits, weights, capacity)

#  Time complexity of this solution is O(n * c) because of the recursive calls are memoized.
#  Here n is the number of items and c is the capacity of the knapsack.
def solve_knapsack_recursive(dp, profits, weights, capacity, index = 0):
    if index >= len(profits) or capacity <= 0:
        return 0

    if dp[index][capacity] != -1:
        return dp[index][capacity]

    p1, p2 = 0, 0
    if weights[index] <= capacity:
        p1 = profits[index] + solve_knapsack_recursive(dp, profits, weights, capacity-weights[index], index+1)

    p2 = solve_knapsack_recursive(dp, profits, weights, capacity, index+1)

    dp[index][capacity] = max(p1, p2)
    return dp[index][capacity]



#  In the recursive solution above we can see that in the recursive calls, profits and weights remain constant.
#  The thing that keeps changing is the capacity and the current index. Now, there are overlapping sub problems that we are solving
#  multiple times repeatedly for the same capacity and current index. In the top down dynamic programming approach, we will store the results
#  of these subproblems in a 2d array so that we can use the memoized results to more easily calculate the final outcome.

def main():
    describe()

    profits = [1, 6, 10, 16]
    weights = [1, 2, 3, 5]
    capacity = 5
    print("Input :")
    print("  Profits : " + str(profits) + " , Weights : " + str(weights) + " , Capacity : " + str(capacity))
    print("  Output : " + str(solve_knapsack(profits, weights, capacity)))
    print("\n")

    profits = [1, 6, 10, 16]
    weights = [1, 2, 3, 5]
    capacity = 6
    print("Input :")
    print("Profits : " + str(profits) + " , Weights : " + str(weights) + " , Capacity : " + str(capacity))
    print("  Output : " + str(solve_knapsack(profits, weights, capacity)))
    print("\n")

    profits = [1, 6, 10, 16]
    weights = [1, 2, 3, 5]
    capacity = 7
    print("Input :")
    print("Profits : " + str(profits) + " , Weights : " + str(weights) + " , Capacity : " + str(capacity))
    print("  Output : " + str(solve_knapsack(profits, weights, capacity)))

main()
