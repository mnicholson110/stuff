from sympy import Matrix

def u(n):
    return 1 - n + n**2 - n**3 + n**4 - n**5 + n**6 - n**7 + n**8 - n**9 + n**10

def V(n):
        return [[x**i for i in range(0,n)] for x in range(1,n+1)]

def BOP(n):
    v = V(n)
    P = []
    for i in v:
        i += [u(i[1])]
    for i in range(n,n*(n+1) + 1,n+1):
        P += [Matrix(v).rref()[0][i]]
    return P

def FIT(n):
    if n == 1: return 1
    P = BOP(n)
    x = n+1
    return sum([P[i]*(x**i) for i in range(len(P))])

print(sum([FIT(i) for i in range(1,11)]))
