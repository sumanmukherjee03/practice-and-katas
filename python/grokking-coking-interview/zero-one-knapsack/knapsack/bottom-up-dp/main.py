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

#  In the bottom up approach to dynamic programming for the knapsack problem
#  we want to find the max profit for every subset of elements and every possible capacity.
#  ie, dp[i][c] represents the max profit for a capacity c calcutaed from first i items.
#  Remember that when we say max profit calculated from first i elements it means that the first i elements were considered, not necessarily included.
#  There are 2 options when calculating dp[i][c]
#       - include the item at index i, ie, dp[i][c] = profit[i] + dp[i-1][c-weight[i]]
#       - exclude the item at index i, ie, dp[i][c] = dp[i-1][c]
#  Our final solution will be max of the 2 options above

#  Time complexity of this solution is O(n * c)
def solve_knapsack(profits, weights, capacity):
    n = len(profits)

    if len(weights) != len(profits) or len(profits) == 0 or capacity <= 0:
        return 0

    #  Create a 2d array with each cell initialized to -1.
    #  The rows represent index of the element upto which the subset of elements is being considered in the solution set.
    #  The columns here represent each capacity broken down by an integer value.
    #  And a cell represents the max profit from considering a subset of elements upto that index and for that max capacity.
    dp = [[0 for x in range(capacity+1)] for y in range(len(profits))]

    #  For capacity 0 populate the first column of all rows with 0 as profit
    for i in range(0, n):
        dp[i][0] = 0

    #  For the 0th item (ie if we only include the 0th item) fill the values in the row with different capacities
    for c in range(0, capacity+1):
        if weights[0] <= c:
            dp[0][c] = profits[0]

    for i in range(1, n):
        for c in range(1, capacity+1):
            p1, p2 = 0, 0
            if weights[i] <= c:
                p1 = profits[i] + dp[i-1][c-weights[i]]
            p2 = dp[i-1][c]
            dp[i][c] = max(p1,p2)

    print_selected_elements(dp, profits, weights, capacity)

    return dp[n-1][capacity]


def print_selected_elements(dp, profits, weights, capacity):
    print("Selected weights are: ", end='')
    n = len(weights)

    #  The bottom right corner cell represents the max profit or the final solution
    total_profit = dp[n-1][capacity]

    #  Iterate from n-1 th index onwards to 0 moving back by 1.
    #  This means we are considering 1 less element at a time as we iterate.
    #  This is gonna tell us if that element was considered or not.
    for i in range(n-1, 0, -1):
        #  If the cumulative profit at this cell is the same as [i-1][capacity] that means this element was not considered in the solution.
        #  But if it is not the same, it indicates that this element was considered
        if total_profit != dp[i-1][capacity]:
            print(str(weights[i]) + " ", end='')
            capacity -= weights[i]
            total_profit -= profits[i]
    if total_profit != 0:
        print(str(weights[0]) + " ", end='')
    print()


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
