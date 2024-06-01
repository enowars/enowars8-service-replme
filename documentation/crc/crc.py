import sympy as sp
from sympy.polys.domains import GF
from itertools import product

x = sp.symbols("x")

# 011xxxxx011xxxxx00000000
#
#      20        16        12         8         4         0
# 0 1 1 k | l m n o | 0 1 1 r | s t u v | 0 0 0 0 | 0 0 0 0
a = (
    x**22
    + x**21
    + sp.symbols("k") * x**20
    + sp.symbols("l") * x**19
    + sp.symbols("m") * x**18
    + sp.symbols("n") * x**17
    + sp.symbols("o") * x**16
    + x**14
    + x**13
    + sp.symbols("r") * x**12
    + sp.symbols("s") * x**11
    + sp.symbols("t") * x**10
    + sp.symbols("u") * x**9
    + sp.symbols("v") * x**8
)

# CRC8-DVB-S2: 0xD5
#
# 8                 0
# 1 | 1 1 0 1 0 1 0 1
p = x**8 + x**7 + x**6 + x**4 + x**2 + 1

# CRC8(aa) = 0xed = CRC8(lc) = CRC(vg) = CRC({e)
#
#
# 1 1 1 0 1 1 0 1
r = x**7 + x**6 + x**5 + x**3 + x**2 + 1

q, rem = sp.div(a, p, domain=GF(2))

equations = sp.Eq(rem, r, domain=GF(2))

solution = sp.solve(
    equations,
    (
        sp.symbols("k"),
        sp.symbols("l"),
        sp.symbols("m"),
        sp.symbols("n"),
        sp.symbols("o"),
        sp.symbols("r"),
        sp.symbols("s"),
        sp.symbols("t"),
        sp.symbols("u"),
        sp.symbols("v"),
    ),
    domain=GF(2)
)

for item in solution:
    print(item)
    free_symbols = set()
    for i in item:
        free_symbols |= i.free_symbols
    possible_values = [0, 1]
    combinations = [
        dict(zip(free_symbols, values))
        for values in product(possible_values, repeat=len(free_symbols))
    ]
    for combo in combinations:
        for eq in item:
            res = eq.subs(combo, domain=GF(2)) % 2
            if res == 0 or res == 1:
                print(combo, "=>", res, eq)
