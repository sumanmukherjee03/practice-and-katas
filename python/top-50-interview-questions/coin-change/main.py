def describe():
    desc = """
Problem : Given an amount of money and a set of possible coins (remember these are coin denominations only, not the actual coins)
          return the minimum number of coins to needed to make that amount. If no such combo exists then return -1.

Example :
    Input : amount = 15, coins = {2,3,7}
    Output : 4
    Explanation : 1*(7) + 2*(3) + 1*(2)

    Input : amount = 34, coins = {5,13}
    Output : -1
    Explanation : No such combination exists

---------------------

    """
    print(desc)

#  This is very similar to the ways to climb stairs problem
#  General pattern is given a target and set of possible denominations try to achieve the target. You can use multiple of the same denominations.
#  General solution to problems like this would be
#       waysToGet(15) = waysToGet(13) + waysToGet(12) + waysToGet(8)
#  So, following a similar approach the solution to this problem will be
#       coinChange(15) = 1 + min(coinChange(13), coinChange(12), coinChange(8))
#  Time complexity of this solution is O(m^n) - where m is the size of the set of the denominations and n is the amount
def coinChange(amount, coins):
    result = coinChangeRecursiveLessConcise(amount, coins)
    return -1 if result == float("inf") else result

def coinChangeRecursiveLessConcise(amount, coins):
    #  If amount is 0 that means you have reached the goal and return 0 to end the recursion
    if amount == 0:
        return 0
    #  Otherwise, consider each coin denomination in the set of coin denominations
    #  For each coin denomination :
    #       If the value of the coin is less than the amount still required to be met as a target
    #           - Find out recursively if there is a way to reach the target by taking this coin denomination
    #       And out of the choices whichever choice yields the minimum number of coins required is considered as the solution
    minNoCoins = [float("inf")]
    for coin in coins:
        if amount - coin >= 0:
            minNoCoins.append(coinChangeRecursiveLessConcise(amount-coin, coins))
    return 1 + min(minNoCoins)

#  Same solution as above. The loop structure is just slightly different.
def coinChangeRecursive(amount, coins):
    #  If amount is 0 that means you have reached the goal and return 0 to end the recursion because there are no more coins to be considered
    if amount == 0:
        return 0
    #  Otherwise, consider each coin denomination in the set of coin denominations
    #  For each coin denomination :
    #       If the value of the coin is less than the amount still required to be met as a target
    #           - Find out recursively if there is a way to reach the target by taking this coin denomination
    #       And out of the choices whichever choice yields the minimum number of coins required is considered as the solution
    minNoCoins = float("inf")
    for coin in coins:
        if amount - coin >= 0:
            minNoCoins = min(minNoCoins, 1 + coinChangeRecursive(amount-coin, coins))
    return minNoCoins








#  The DP solution follows the path that to reach the final target amount, we try and find ways to reach each amount from 0 -> amount using the coin denominations.
#  dp[i] = 1 + min(dp[i-coin1], dp[i-coin2], ....)
#  ie minimum of the choices for taking various coin denominations
#  Time complexity of the DP solution is O(n)
def coinChangeDP(amount, coins):
    nbCoinsArr = [float("inf")] * (amount+1)
    nbCoinsArr[0] = 0
    for i in range(1, amount+1):
        minCoins = float("inf")
        for coin in coins:
            if (i-coin) >= 0:
                minCoins = min(minCoins, 1 + nbCoinsArr[i-coin])
        nbCoinsArr[i] = minCoins
    return -1 if nbCoinsArr[amount] == float("inf") else nbCoinsArr[amount]








def main():
    describe()

    amount = 15
    coins = [2,3,7]
    print("Input : " + str(amount) + ", " + str(coins))
    print("Output : " + str(coinChangeDP(amount, coins)))
    print("\n\n")

    amount = 34
    coins = [5,13]
    print("Input : " + str(amount) + ", " + str(coins))
    print("Output : " + str(coinChangeDP(amount, coins)))
    print("\n\n")

main()
