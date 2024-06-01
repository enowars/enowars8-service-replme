import random
from collections import namedtuple

def str_to_long(s, reverse_bytes=False):
    ret = 0
    for c in s:
        ret = (ret << 8)
        c = ord(c)
        if reverse_bytes:
            c = int('{:08b}'.format(c)[::-1], 2)
        ret = ret ^^ c

    return ret

def long_to_str(l, reverse_bytes=False):
    s = []
    while (l > 0):
        n = l&0xff
        if reverse_bytes:
            n = int('{:08b}'.format(n)[::-1], 2)
        s.append(n)
        l = l >> 8
    s.reverse()
    return bytes(s)

def alphanumeric_to_long(s, reverse_bytes=False):
    ret = 0
    for c in s:
        ret = (ret << 8)
        v = ord(c)
        if reverse_bytes:
            v = int('{:08b}'.format(v)[::-1], 2)
        if v == 0b110000:
            ret = ret ^^ 0b1100000
        elif v >= 0b110010 and v <= 0b110101:
            ret = ret ^^ (v+74)
        else:
            ret = ret ^^ v
    return ret

def long_to_alphanumeric(l, reverse_bytes=False):
    s = []
    while (l > 0):
        c = l&0xff
        if reverse_bytes:
            c = int('{:08b}'.format(c)[::-1], 2)
        if c == 0b1100000:
            s.append(0b110000)
        elif c >= 0b1111011 and c <= 0b01111111:
            s.append(c-74)
        else:
            s.append(c)
        l = l >> 8
    s.reverse()
    return bytes(s).decode()

def crcN(p, s):
    d = len(bin(p)) - 3
    mask = (2**d) - 1
    creg = 0
    for j in range(len(s)):
        c = ord(s[j])
        for _ in range(8):
            if ((creg ^^ (c << (d-8))) & (1<<(d-1))) > 0:
                creg = ((creg << 1) ^^ p) & mask
            else:
                creg = (creg << 1) & mask
            c = c << 1
    return creg

Context = namedtuple('Context', 'degree P x bs p b pb')

def get_random_prime_polynom(degree):
    P.<x> = PolynomialRing(Zmod(2), 'x')
    while(True):
        p = P.random_element(degree);
        if str(factor(p)) == str(p):
            return int(''.join(map(str, p.polynomial(x).coefficients(sparse=False)[::-1])), 2)

def prepare(polynom, input_len):
    # prepare ring and vars
    degree = len(bin(polynom)) - 3

    __v = ['x']

    for i in range(input_len*8):
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
    for i in range(input_len*8):
        b += bs[i] * x^i

    pb = p*b

    return Context(degree, P, x, bs, p, b, pb)

def calculate_kernel(ctx, mask):
    # prepare underdetermined system for
    x = ctx.x
    A = []
    for i in range(ctx.pb.degree()-1, ctx.degree, -8):
        for j, v in enumerate(mask):
            if v == 0:
                continue
            row = [0]*len(ctx.bs)
            for c, _b in ctx.pb.coefficient(x^(i-j)):
                row[ctx.bs.index(_b)] = 1
            A.append(row)

    for i in range(ctx.degree, 0, -1):
        row = [0]*len(ctx.bs)
        for c, _b in ctx.pb.coefficient(x^i):
            row[ctx.bs.index(_b)] = 1
        A.append(row)

    row = [0]*len(ctx.bs)
    for c, _b in ctx.pb(x=0):
        row[ctx.bs.index(_b)] = 1

    A.append(row)

    # calculate kernel span Ab=0
    A = Matrix(GF(2), A)
    k = A.right_kernel().basis()
    return k

def get_random_delta(ctx, k):
    x = ctx.x
    _c = random.randrange(1, 2**(len(k)))
    _v = [0]*len(k)
    for i in range(len(k)):
        _v[i] = (_c >> i) & 1
    _d = sum(c*v for c, v in zip(_v, k))
    d = 0
    for i, bi in enumerate(_d):
        d += bi*x^i
    d = d*ctx.p
    __c = d.univariate_polynomial().coefficients(sparse=False)
    return ''.join(map(str, __c[ctx.p.degree():][::-1]))

def calculate_pre_image(p, s, reverse_in=False):
    print("Polynom     :", hex(p))
    print("Degree      :", len(bin(p)) - 3)
    print("Input       :", s)
    # print("CRCN(Input) :", hex(crcN(p, s)))
    ctx = prepare(p, len(s))
    mask = [0,0,0,0,0,1,1,1] if reverse_in else [1,1,1,0,0,0,0,0]
    k = calculate_kernel(ctx, mask)
    print("#Span       :", len(k))
    print("#Kernel     :", 2**(len(k)))
    while True:
        d = int(get_random_delta(ctx, k), 2)
        s2 = str_to_long(s, reverse_in) ^^ d
        s2raw = long_to_str(s2, reverse_in)
        s2an = long_to_alphanumeric(s2, reverse_in)
        print("raw pre-img :", s2raw)
        print("2nd pre-img :", s2an)
        if s2an[0].isalpha():
            break
    print()

# calculate_pre_image(0x104C11DB7, "aaaaaaaaaaaaaaaaaaaaaaaa", True)
calculate_pre_image(0x104C11DB7, "aaaaaaaaaaaaaaaaaaaaaaaa")
# calculate_pre_image(get_random_prime_polynom(187), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
