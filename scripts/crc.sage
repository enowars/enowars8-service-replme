# basic prepare ring and vars
polynom = 0b111010101
degree = len(bin(polynom)) - 3

__v = ['x']

for i in range(degree):
    __v.append(f"b{i}")

P = PolynomialRing(Zmod(2), __v)

x = P.gens()[0]
bs = P.gens()[1:]

# create polynom p
p = 0
for i in range(degree+1):
    p += ((polynom >> i) & 1)*x^i

# create coefficient b (a=bp+r)
b = 0
for i in range(degree):
    b += bs[i] * x^i

pb = p*b

# prepare underdetermined system for
A = []
for i in range(15, 0, -8):
    for j in range(3):
        row = [0]*len(bs)
        for c, _b in pb.coefficient(x^(i-j)):
            row[bs.index(_b)] = 1
        A.append(row)

# calculate kernel span Ab=0
A = Matrix(GF(2), A)
k = A.right_kernel().basis()

# calculate all elements in kernel
def all_combinations(vectors):
    from itertools import product
    n = len(vectors)
    for comb in product([0, 1], repeat=n):
        yield sum(c*v for c, v in zip(comb, vectors))

# calculate delta for each element in kernel d=ep
solutions = []
for e in all_combinations(k):
    solution = 0
    for i, bi in enumerate(e):
        solution += bi*x^i
    solutions.append(p*solution)

# delta to binary
for s in solutions:
    __c = s.univariate_polynomial().coefficients(sparse=False)
    if len(__c) > 0:
        print(''.join(map(str, __c[::-1])))

# b2 = 1*x^4 + 1*x^3 + 1*x^2 + 1*x
# b3 = 0*x^4 + 1*x^3 + 0*x^2 + 1*x
# b4 = 1*x^4 + 0*x^3 + 1*x^2 + 0*x
# b5 = 0*x^4 + 0*x^3 + 0*x^2 + 0*x

# print("===")
# print(p*b2)
# print("===")
# print(p*b3)
# print("===")
# print(p*b4)
# print("===")
# print(p*b5)
# print("===")

# b7 = 0
# b6 = 0
# b5 = 0
# b0 = 0
#
# b2 = b4
# b1 = b3
#

#          12  1110 9 8   7 6 5 4   3 2 1 0
#           k   l m n o         r   s t u v
# aa  0 1 1 0 | 0 0 0 1 | 0 1 1 0 | 0 0 0 1
#  +  0 0 0 1 | 0 1 1 1 | 0 0 0 0 | 0 1 1 0
# ??  0 1 1 1 | 0 1 1 0 | 0 1 1 0 | 0 1 1 1
# vg  0 1 1 1 | 0 1 1 0 | 0 1 1 0 | 0 1 1 1


# aa  0 1 1 0 | 0 0 0 1 | 0 1 1 0 | 0 0 0 1
# ??  0 1 1 0 | 1 1 0 0 | 0 1 1 0 | 0 0 1 1
# lc  0 1 1 0 | 1 1 0 0 | 0 1 1 0 | 0 0 1 1


# aa  0 1 1 0 | 0 0 0 1 | 0 1 1 0 | 0 0 0 1
# ??  0 1 1 1 | 1 0 1 1 | 0 1 1 0 | 0 1 0 1
# {e  0 1 1 1 | 1 0 1 1 | 0 1 1 0 | 0 1 0 1
#
#
# 0 = b7
# 0 = b6 + b7
# 0 = b5 + b6 + b7
# 0 = b0 + b1 + b3 + b5 + b7
# 0 = b0 + b2 + b4 + b6
# 0 = b1 + b3 + b5
#
#
# 1 0 0 0 0 0 0 0   b7   0
# 1 1 0 0 0 0 0 0   b6   0
# 1 1 1 0 0 0 0 0   b5   0
# 1 0 1 0 1 0 1 1 * b4 = 0
# 0 1 0 1 0 1 0 1   b3   0
# 0 0 1 0 1 0 1 0   b2   0
#                   b1
#                   b0
#
#
