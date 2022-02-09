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
    return solve_knapsack_recursive(profits, weights, capacity)

def solve_knapsack_recursive(profits, weights, capacity, index = 0):
    if index >= len(profits) or capacity <= 0:
        return 0

    p1, p2 = 0, 0
    if weights[index] <= capacity:
        p1 = profits[index] + solve_knapsack_recursive(profits, weights, capacity-weights[index], index+1)

    p2 = solve_knapsack_recursive(profits, weights, capacity, index+1)
    return max(p1, p2)


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
