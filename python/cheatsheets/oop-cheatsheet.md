```
class Dog(object):
  living = True
  age = 0

  def __init__(self, name, sex):
    self.name = name
    self.sex = sex

  def bark(self):
    print(f"{self.name} barked!")


luke = Dog('Luke', 'Male')
print(f"Name: {luke.name}")
print(f"Sex: {luke.sex}")
print(f"Living: {luke.living}")
print(f"Age: {luke.age}")
luke.bark()
```

```
class Vehicle(object):
  def honk(self):
    print("honk honk")

class Car(Vehicle):
  def accelerate:
    print("vrooom")

honda = Car()
honda.honk()
honda.accelerate()
```


Python class magic methods
```
__init__(self, other)  # Overrides object creation method
__repr__(self)         # Overrides object string representation
__add__(self, other)   # Overrides + operator
__sub__(self, other)   # Overrides - operator
__mul__(self, other)   # Overrides * operator
__floordiv__(self, other)   # Overrides // operator
__truediv__(self, other)   # Overrides / operator
__mod__(self, other)   # Overrides % operator
__lt__(self, other)   # Overrides < comparison operator
__le__(self, other)   # Overrides <= comparison operator
__eq__(self, other)   # Overrides == comparison operator
__ne__(self, other)   # Overrides != comparison operator
__gt__(self, other)   # Overrides > comparison operator
__ge__(self, other)   # Overrides >= comparison operator
__call__(self[, args...])   # Overrides () operator
__int__(self)   # Overrides int() method
__float__(self)   # Overrides float() method
__abs__(self)   # Overrides abs() method
__len__(self)   # Overrides len() method
__contains__(self, item)   # Overrides contains() method
```
