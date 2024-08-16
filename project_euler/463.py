memo = {0:0,1:1,2:2,3:5}

def S(n):
    sol = memo.get(n)
    if sol is not None:
        return sol
    
    q, r = divmod(n,4)
    if r == 0:
        sol = S(q*2)*6 - S(q)*5 - S(q - 1)*3 - 1
    elif r == 1:
        sol = S(q*2 + 1)*2 + S(q*2)*4 - S(q)*6 - S(q - 1)*2 - 1
    elif r == 2:
        sol = S(q*2 + 1)*3 + S(q*2)*3 - S(q)*6 - S(q - 1)*2 - 1
    else:
        sol = S(q*2 + 1)*6 - S(q)*8 - 1

    memo[n] = sol
    return sol
    
    
print(S(3**37)%10**9)
