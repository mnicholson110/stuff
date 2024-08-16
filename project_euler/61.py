from itertools import permutations

def triangle(i):
  return i*(i+1)/2
def square(i):
  return i*i
def pentagonal(i):
  return i*(3*i - 1)/2
def hexagonal(i):
  return i*(2*i - 1)
def heptagonal(i):
  return i*(5*i - 3)/2
def octagonal(i):
  return i*(3*i - 2)

checks = [[triangle(i) for i in range(45,141)],
[square(i) for i in range(32,100)],
[pentagonal(i) for i in range(25,82)],
[hexagonal(i) for i in range(23,71)],
[heptagonal(i) for i in range(21,64)],
[octagonal(i) for i in range(19,59)]]

def func61((a,b,c,d,e,f)):
  for n1 in checks[a]:
    for n2 in checks[b]:
      if n2 % 100 == n1/100:
        for n3 in checks[c]:
          if n3 % 100 == n2/100:
            for n4 in checks[d]:
              if n4 % 100 == n3/100:
                for n5 in checks[e]:
                  if n5 % 100 == n4/100:
                    for n6 in checks[f]:
                      if n6 % 100 == n5/100:
                        if n1 % 100 == n6/100:
                          print (n1,n2,n3,n4,n5,n6)
                          print sum((n1,n2,n3,n4,n5,n6))

for p in [i for i in permutations([1,2,3,4,5])]:
  func61((0,p[0],p[1],p[2],p[3],p[4]))
