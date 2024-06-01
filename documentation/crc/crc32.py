import sympy as sp
from sympy.polys.domains import GF
from itertools import product
import time

x = sp.symbols("x")

a = 0
s = []

for i in range(32, (32 + 8 * 7), 8):
    s4 = sp.symbols("e" + str(i))
    s3 = sp.symbols("d" + str(i))
    s2 = sp.symbols("c" + str(i))
    s1 = sp.symbols("b" + str(i))
    s0 = sp.symbols("a" + str(i))
    a += (
        (x ** (i + 6))
        + (x ** (i + 5))
        + s4 * x ** (i + 4)
        + s3 * x ** (i + 3)
        + s2 * x ** (i + 2)
        + s1 * x ** (i + 1)
        + s0 * x ** (i)
    )
    s.extend([s0, s1, s2, s3, s4])

print("a:", a)

# for i in range(32, (32 + 8 * 3), 8):
#     s4 = sp.symbols("e" + str(i))
#     s3 = sp.symbols("d" + str(i))
#     s2 = sp.symbols("c" + str(i))
#     s1 = sp.symbols("b" + str(i))
#     s0 = sp.symbols("a" + str(i))
#     a += (
#         s4 * x ** (i + 7)
#         + s3 * x ** (i + 6)
#         + s2 * x ** (i + 5)
#         + s1 * x ** (i + 4)
#         + s0 * x ** (i + 3)
#         + (x ** (i + 2))
#         + (x ** (i + 1))
#     )
#     s.extend([s0, s1, s2, s3, s4])
#
# for i in range((32 + 8 * 3), (32 + 8 * 7), 8):
#     s4 = sp.symbols("e" + str(i))
#     s3 = sp.symbols("d" + str(i))
#     s2 = sp.symbols("c" + str(i))
#     s1 = sp.symbols("b" + str(i))
#     s0 = sp.symbols("a" + str(i))
#     a += (
#         (s4 + 1) * x ** (i + 7)
#         + (s3 + 1) * x ** (i + 6)
#         + (s2 + 1) * x ** (i + 5)
#         + (s1 + 1) * x ** (i + 4)
#         + (s0 + 1) * x ** (i + 3)
#         + (x ** (i))
#     )
#     s.extend([s0, s1, s2, s3, s4])

poly = 0x104C11DB7
p = 0

for i in range(33):
    if (poly >> i) & 1 == 1:
        p += 1 * x**i

print("p:", p)

# rest = 0xCF8101E4
rest = 0xB70ED28
r = 0

for i in range(32):
    if (rest >> i) & 1 == 1:
        r += 1 * x**i

print("r:", r)

start = time.time()
q, rem = sp.div(a, p, domain=GF(2))
end = time.time()
print(end - start)


start = time.time()
equations = sp.Eq(rem, r)
end = time.time()
print(end - start)

s.reverse()

start = time.time()
solution = sp.solve(equations, s, domain=GF(2))
end = time.time()
print(end - start)

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
            res = eq.subs(combo) % 2
            print(combo, "=>", res, "(", eq, ")")
