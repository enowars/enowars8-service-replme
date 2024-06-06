from random import randrange

NOISE = [
    (
        [  # pseudo-user #1
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
        [  # validation #1
            ("gcc -o main main.c && echo OK\n", ".*\nOK.*"),
            ("./main\n", ".*\nHello World!.*"),
        ],
    ),
    (
        [  # pseudo-user #2
            ("mkdir docs && echo OK\n", ".*\nOK.*"),
            ("cd docs && echo OK\n", ".*\nOK.*"),
            ("ls -lah && echo OK\n", ".*\nOK.*"),
            ('echo "And I say hey" && echo OK\n', ".*\nOK.*"),
            ('echo "What\'s going on" && echo OK\n', ".*\nOK.*"),
            ("touch 4-non-blondes.txt && echo OK\n", ".*\nOK.*"),
            ('echo "HEY, YEAH, YEAH" >> 4-non-blondes.txt && echo OK\n', ".*\nOK.*"),
        ],
        [  # validation #2
            ("cd docs && echo OK\n", ".*\nOK.*"),
            ("cat 4-non-blondes.txt\n", ".*\nHEY, YEAH, YEAH.*"),
        ],
    ),
    (
        [  # pseudo-user #3
            ("touch notes.txt && echo OK\n", ".*\nOK.*"),
            ('echo "To do list:" > notes.txt && echo OK\n', ".*\nOK.*"),
            ('echo "- Buy groceries" >> notes.txt && echo OK\n', ".*\nOK.*"),
            ('echo "- Call mom" >> notes.txt && echo OK\n', ".*\nOK.*"),
            ('echo "- Finish project" >> notes.txt && echo OK\n', ".*\nOK.*"),
            ("cat notes.txt && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #3
            (
                "cat notes.txt\n",
                ".*\nTo do list:.*\n- Buy groceries.*\n- Call mom.*\n- Finish project.*",
            ),
        ],
    ),
    (
        [  # pseudo-user #4
            ("mkdir project && echo OK\n", ".*\nOK.*"),
            ("cd project && echo OK\n", ".*\nOK.*"),
            ('git config --global user.email "1337@pwn.y" && echo OK\n', ".*\nOK.*"),
            ('git config --global user.name "Eve" && echo OK\n', ".*\nOK.*"),
            ("git init && echo OK\n", ".*\nOK.*"),
            ('echo "# My Project" > README.md && echo OK\n', ".*\nOK.*"),
            ("git add README.md && echo OK\n", ".*\nOK.*"),
            ('git commit -m "Initial commit" && echo OK\n', ".*\nOK.*"),
        ],
        [  # validation #4
            ("cd project && echo OK\n", ".*\nOK.*"),
            ("git --no-pager log --oneline\n", ".*Initial commit.*"),
        ],
    ),
    (
        [  # pseudo-user #5
            ("mkdir mydir && echo OK\n", ".*\nOK.*"),
            ("cd mydir && echo OK\n", ".*\nOK.*"),
            ("touch file1.txt && echo OK\n", ".*\nOK.*"),
            ("touch file2.txt && echo OK\n", ".*\nOK.*"),
            ("ls -l && echo OK\n", ".*\nOK.*"),
            ("mv file1.txt file3.txt && echo OK\n", ".*\nOK.*"),
            ("ls -l && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #5
            ("cd mydir && echo OK\n", ".*\nOK.*"),
            ("ls -l\n", ".*file2.txt.*file3.txt.*"),
        ],
    ),
    (
        [  # pseudo-user #6
            ("echo 'Hello, World!' > greeting.txt && echo OK\n", ".*\nOK.*"),
            ("cat greeting.txt && echo OK\n", ".*\nOK.*"),
            ("cp greeting.txt copy_of_greeting.txt && echo OK\n", ".*\nOK.*"),
            ("cat copy_of_greeting.txt && echo OK\n", ".*\nOK.*"),
            ("rm greeting.txt && echo OK\n", ".*\nOK.*"),
            ("ls && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #6
            ("cat copy_of_greeting.txt\n", ".*Hello, World!.*"),
        ],
    ),
    (
        [  # pseudo-user #7
            ("echo 'User 7 logging in' > log.txt && echo OK\n", ".*\nOK.*"),
            ("echo 'Performing task 1' >> log.txt && echo OK\n", ".*\nOK.*"),
            ("echo 'Performing task 2' >> log.txt && echo OK\n", ".*\nOK.*"),
            ("echo 'User 7 logging out' >> log.txt && echo OK\n", ".*\nOK.*"),
            ("cat log.txt && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #7
            (
                "cat log.txt\n",
                ".*User 7 logging in.*Performing task 1.*Performing task 2.*User 7 logging out.*",
            ),
        ],
    ),
    (
        [  # pseudo-user #8
            ("mkdir music && echo OK\n", ".*\nOK.*"),
            ("cd music && echo OK\n", ".*\nOK.*"),
            ("touch song1.mp3 song2.mp3 && echo OK\n", ".*\nOK.*"),
            ("ls && echo OK\n", ".*\nOK.*"),
            ("mv song1.mp3 favorite_song.mp3 && echo OK\n", ".*\nOK.*"),
            ("ls && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #8
            ("cd music && echo OK\n", ".*\nOK.*"),
            ("ls\n", ".*favorite_song.mp3.*song2.mp3.*"),
        ],
    ),
    (
        [  # pseudo-user #9
            ("touch data.csv && echo OK\n", ".*\nOK.*"),
            ('echo "name,age" > data.csv && echo OK\n', ".*\nOK.*"),
            ('echo "Alice,30" >> data.csv && echo OK\n', ".*\nOK.*"),
            ('echo "Bob,25" >> data.csv && echo OK\n', ".*\nOK.*"),
            ("cat data.csv && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #9
            ("cat data.csv\n", ".*name,age.*Alice,30.*Bob,25.*"),
        ],
    ),
    (
        [  # pseudo-user #10
            ("mkdir scripts && echo OK\n", ".*\nOK.*"),
            ("cd scripts && echo OK\n", ".*\nOK.*"),
            ("touch script.sh && echo OK\n", ".*\nOK.*"),
            ("chmod +x script.sh && echo OK\n", ".*\nOK.*"),
            ("echo '#!/bin/sh' > script.sh && echo OK\n", ".*\nOK.*"),
            ('echo "echo Hello, Script" >> script.sh && echo OK\n', ".*\nOK.*"),
            ("./script.sh && echo OK\n", ".*\nHello, Script.*"),
        ],
        [  # validation #10
            ("cd scripts && echo OK\n", ".*\nOK.*"),
            ("./script.sh\n", ".*Hello, Script.*"),
        ],
    ),
    (
        [  # pseudo-user #11
            ("mkdir photos && echo OK\n", ".*\nOK.*"),
            ("cd photos && echo OK\n", ".*\nOK.*"),
            ("touch img1.jpg img2.jpg img3.jpg && echo OK\n", ".*\nOK.*"),
            ("ls && echo OK\n", ".*\nOK.*"),
            ("rm img2.jpg && echo OK\n", ".*\nOK.*"),
            ("ls && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #11
            ("cd photos && echo OK\n", ".*\nOK.*"),
            ("ls\n", ".*img1.jpg.*img3.jpg.*"),
        ],
    ),
    (
        [  # pseudo-user #12
            ("touch todo.txt && echo OK\n", ".*\nOK.*"),
            ('echo "1. Finish homework" > todo.txt && echo OK\n', ".*\nOK.*"),
            ('echo "2. Clean room" >> todo.txt && echo OK\n', ".*\nOK.*"),
            ('echo "3. Exercise" >> todo.txt && echo OK\n', ".*\nOK.*"),
            ("cat todo.txt && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #12
            (
                "cat todo.txt\n",
                ".*1. Finish homework.*2. Clean room.*3. Exercise.*",
            ),
        ],
    ),
    (
        [  # pseudo-user #13
            ("mkdir workspace && echo OK\n", ".*\nOK.*"),
            ("cd workspace && echo OK\n", ".*\nOK.*"),
            ("touch project1.txt project2.txt && echo OK\n", ".*\nOK.*"),
            ("ls && echo OK\n", ".*\nOK.*"),
            ("rm project1.txt && echo OK\n", ".*\nOK.*"),
            ("ls && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #13
            ("cd workspace && echo OK\n", ".*\nOK.*"),
            ("ls\n", ".*project2.txt.*"),
        ],
    ),
    (
        [  # pseudo-user #14
            ("echo 'Hello from User 42' > message.txt && echo OK\n", ".*\nOK.*"),
            ("cat message.txt && echo OK\n", ".*\nOK.*"),
            ("mv message.txt greetings.txt && echo OK\n", ".*\nOK.*"),
            ("cat greetings.txt && echo OK\n", ".*\nOK.*"),
        ],
        [  # validation #14
            ("cat greetings.txt\n", ".*Hello from User 42.*"),
        ],
    ),
]


def get_noise(i: int):
    return NOISE[i]


def get_random_noise():
    i = randrange(len(NOISE))
    return (i, NOISE[i])
