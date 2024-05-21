from random import randrange

NOISE = [
    (
        [
            ('echo "#include <stdio.h>" > main.c && echo OK\n', ".*\nOK.*"),
            ('echo "int main() {" >> main.c && echo OK\n', ".*\nOK.*"),
            (
                "echo '    printf(\"Hello World!\\\\n\");' >> main.c && echo OK\n",
                ".*\nOK.*",
            ),
            ('echo "    return 0;" >> main.c && echo OK\n', ".*\nOK.*"),
            ('echo "}" >> main.c && echo OK\n', ".*\nOK.*"),
            ("gcc -o main main.c && echo OK\n", ".*\nOK.*"),
            ("./main\n", ".*\nHello World!.*"),
        ],
        [
            ("gcc -o main main.c && echo OK\n", ".*\nOK.*"),
            ("./main\n", ".*\nHello World!.*"),
        ],
    )
]


def get_noise(i: int):
    return NOISE[i]


def get_random_noise():
    i = randrange(len(NOISE))
    return (i, NOISE[i])
