### 0/1 knapsack (Dynamic programming)

#### Description
Given the weights and profits of N items, we are asked to put these items in a knapsack with a capacity C.
The goal is to get the maximum profit out of the knapsack items. Each item can only be selected once, as we dont have multiple quantities of any item.

Merry wants to carry some fruits in the knapsack to get maximum profit. Here are the weights and profits of the fruits:

Items: { Apple, Orange, Banana, Melon }
Weights: { 2, 3, 1, 4 }
Profits: { 4, 5, 3, 7 }
Knapsack capacity: 5

Combinations of fruits in the knapsack, such that their total weight is not more than 5:

Apple + Orange (total weight 5) => 9 profit
Apple + Banana (total weight 3) => 7 profit
Orange + Banana (total weight 4) => 8 profit
Banana + Melon (total weight 5) => 10 profit

This shows that Banana + Melon is the best combination as it gives us the maximum profit, and the total weight does not exceed the capacity.



In this pattern there are usually recursive solutions but they can be improved with
techniques like memoization and bottom-up dynamic programming.
