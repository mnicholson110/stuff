from math import floor

def B(R):
    return 0.5*(8*(R**2)+1)**0.5 + 0.5*(2*R+1)
    #x_2 = 0.5*(2*R+1) - 0.5*(8*(R**2)+1)**0.5

    #return [x_1, x_2]

def possibleR(n):
    l = []
    for i in range(292850000001,n,7):
        if floor((8*i**2 + 1)**0.5) == (8*i**2 + 1)**0.5:
            l += [i]
    for i in range(292850000002,n,7):
        if floor((8*i**2 + 1)**0.5) == (8*i**2 + 1)**0.5:
            l += [i]
    for i in range(292850000003,n,7):
        if floor((8*i**2 + 1)**0.5) == (8*i**2 + 1)**0.5:
            l += [i]
    l.sort()
    return l

def P(b,r):
    return (b/(b+r))*((b-1)/(b+r-1))

for (R,T) in [(x, B(x) + x) for x in possibleR(292893200000 + 10**7) if  B(x) == int(B(x))]:
    if T > 10**12:
        print(T-R,R,T,P(T-R,R))
