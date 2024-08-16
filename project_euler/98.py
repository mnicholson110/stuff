def letters(w):
    let = []
    for l in w:
        let += [l]
    let.sort()
    return let

def getDigits(n):
    p = []
    while n:
        p += [n % 10]
        n //= 10
    p.sort()
    return p


def squaresOfLen(n):
    return [i**2 for i in range(31623) if len(getDigits(i**2)) == n]




