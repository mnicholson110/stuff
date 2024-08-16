from math import floor, e

def listSum(*arg):
    s = []
    for a in arg[0]:
        s += a
    return s

def contfrace(n):
    
    return [2] + listSum([[1,k,1] for k in range(2,n,2)])

