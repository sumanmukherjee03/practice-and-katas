#  Problem : Given a staircase of n steps and a set of possible steps that we can climb at a time
#  named possibleSteps, create a function that returns the number of ways a person can take to reach the
#  top of the staircase.
#  Ex : Input : n = 5, possibleSteps = {1,2}
#  Output : 8
#  Explanation : Possible ways are - 11111, 1112, 1121, 1211, 2111, 122, 212, 221


#  Think recursively with a few examples and see if there is any pattern to the solution
#  When n=0 , possibleSteps = {1,2} - Output = 1
#  When n = 1, possibleSteps = {1,2} - Output = 1
#  When n = 2, possibleSteps = {1,2} - Output = 2
#  When n = 3, possibleSteps = {1,2} - Output = 3
#  When n = 4, possibleSteps = {1,2} - Output = 5
#  When n = 5, possibleSteps = {1,2} - Output = 8
#  The output is a fibonacci sequence : 1,1,2,3,5,8
#  The recursive relation here is f(n) = f(n-1) + f(n-2)
#
#  This is because to climb any set of steps n, if the solution is f(n),
#  then we can find f(n) by finding the solution to climb n-1 steps + 1 step OR n-2 steps and 2 steps,
#  ie to reach the top we must be either 1 step away from it or 2 steps away from it.
#  That's why f(n) = f(n-1) + f(n-2)
#
#  Similarly when possibleSteps = {2,3,4}, then to reach any step n, we must be either 2 steps
#  away from it or 3 steps away from it or 4 steps away from it.
#  f(n) = f(n-2) + f(n-3) + f(n-4)

def waysToClimb(n, possibleSteps):
    if n == 0:
        return 1
    noWays = 0
    for steps in possibleSteps:
        if n-steps > 0:
            noWays += waysToClimb(n-steps, possibleSteps)
    return noWays
