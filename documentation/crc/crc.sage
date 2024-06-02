import random
import json
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

def gen_crc_table(p, reverse_in=False):
    degree = len(bin(p)) - 3
    table = []
    mask = 2**degree-1
    for i in range(256):
        c = i
        if reverse_in:
            c = int('{:08b}'.format(i)[::-1], 2)
        c = c << degree - 8
        m = 1 << (degree - 1)
        for _ in range(8):
            if (c & m) != 0:
                c = (c << 1) ^^ p
            else:
                c = c << 1
        c = c & mask
        if reverse_in:
            c = int(bin(c^^(1 << degree))[2::][::-1][:-1], 2)
        table.append(c)
    return table

def crcT(degree, table, s):
    r = 0
    mask = 2**degree-1
    for j in range(len(s)):
        c = ord(s[j])
        r = table[((r >> (degree - 8)) ^^ c) & 0xff] ^^ (r << 8)
        r = r & mask
    return r


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
    checker_file = open('checker_data.json', 'w')
    service_file = open('service_data.json', 'w')
    checker_data = {
        'polynom': hex(p)[2:]
    }
    service_data = {}
    degree = len(bin(p)) - 3
    print("Polynom     :", hex(p))
    print("Degree      :", degree)
    print("Input       :", s)
    # print("CRCN(Input) :", hex(crcN(p, s)))
    ctx = prepare(p, len(s))
    mask = [0,0,0,0,0,1,1,1] if reverse_in else [1,1,1,0,0,0,0,0]
    k = calculate_kernel(ctx, mask)
    print("#Span       :", len(k))
    print("#Kernel     :", 2**(len(k)))
    deltas = []
    for i in range(100):
        d = int(get_random_delta(ctx, k), 2)
        deltas.append(hex(d)[2:])
    checker_data['deltas'] = deltas
    while True:
        d = int(get_random_delta(ctx, k), 2)
        s2 = str_to_long(s, reverse_in) ^^ d
        s2raw = long_to_str(s2, reverse_in)
        s2an = long_to_alphanumeric(s2, reverse_in)
        print("raw pre-img :", s2raw)
        print("2nd pre-img :", s2an)
        if s2an[0].isalpha():
            break
    table = gen_crc_table(p, reverse_in)
    table = list(map(lambda x: hex(x)[2:], table))
    checker_data['table'] = table
    service_data['table'] = table
    checker_file.write(json.dumps(checker_data, indent=2))
    service_file.write(json.dumps(service_data, indent=2))
    checker_file.close()
    service_file.close()
    print()

# calculate_pre_image(get_random_prime_polynom(251), "a"*60)
calculate_pre_image(0x8561cc4ee956c6503c5da0ffacb20feabb3eb142e7645e7ff1a2067fd8e1cfb, "a"*60)
