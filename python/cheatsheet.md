## Python cheatsheet


### operators
```
** - exponent
// - integer division
```

### string
```
'Alice' 'Bob' # AliceBob
'a' * 5 # aaaaa
'what ever the string is, here is the interpolation {}'.format(variable)
len('hello')

"""
I am a docstring
"""
```

### userinput
```
print('What is your name?')
name = input()
print('Good to meet you {}'.format(name))
```

### scope
If you need to modify a global variable from within a function, use the global statement

```
eggs = 'global'
def spam():
  global eggs
  eggs = 'spam'

spam()
print(eggs) # spam
```


### datastructures

```
mylist = otherlist[:] # slicing the entire list performs a copy of the list
```


### miscellaneous
```
values = [0,0,0,0,0]
res = any(values) # False

values = [None, [], "", 0]
res = any(values) # False

str(21)
int(7.7)

# prints even numbers from 0 to 8
for i in range(0, 10, 2):
  print(i)

# print 5 to 0 backwards
for i in range(5, -1, -1):
  print(i)

import random, sys, os, math
random.randint(1, 10)
sys.exit()

print('Hello', end='')
print('world')
# Helloworld

print('dog', 'rabbit', 'cat', sep=',') # dog,rabbit,cat
```


