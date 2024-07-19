from random import randrange

from util import ShellCommandChain, sh, shchain

NOISE: list[ShellCommandChain] = [
    shchain(
        cmds=[
            sh('echo "2b or not 2b" > note.txt').default(),
        ],
        validations=[
            sh("cat note.txt").expect(".*\n2b or not 2b.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "#include <stdio.h>" > main.c').default(),
            sh('echo "int main() {" >> main.c').default(),
            sh("echo '    printf(\"Hello World!\\\\n\");' >> main.c").default(),
            sh('echo "    return 0;" >> main.c').default(),
            sh('echo "}" >> main.c').default(),
            sh("gcc -o main main.c").default(),
            sh("./main").expect(".*\nHello World!.*").default(),
        ],
        validations=[
            sh("gcc -o main main.c").default(),
            sh("./main").expect(".*\nHello World!.*").default(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir docs").default(),
            sh("cd docs").default(),
            sh("ls -lah").default(),
            sh('echo "And I say hey"').default(),
            sh('echo "What\'s going on"').default(),
            sh("touch 4-non-blondes.txt").default(),
            sh('echo "HEY, YEAH, YEAH" >> 4-non-blondes.txt').default(),
        ],
        validations=[
            sh("cd docs").default(),
            sh("cat 4-non-blondes.txt").expect(".*\nHEY, YEAH, YEAH.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("touch notes.txt").default(),
            sh('echo "To do list:" > notes.txt').default(),
            sh('echo "- Buy groceries" >> notes.txt').default(),
            sh('echo "- Call mom" >> notes.txt').default(),
            sh('echo "- Finish project" >> notes.txt').default(),
            sh("cat notes.txt").default(),
        ],
        validations=[
            sh("cat notes.txt")
            .expect(
                ".*\nTo do list:.*\n- Buy groceries.*\n- Call mom.*\n- Finish project.*",
            )
            .errext()
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir project").default(),
            sh("cd project").default(),
            sh('git config --global user.email "1337@pwn.y"').default(),
            sh('git config --global user.name "Eve"').default(),
            sh("git init").default(),
            sh('echo "# My Project" > README.md').default(),
            sh("git add README.md").default(),
            sh('git commit -m "Initial commit"').default(),
        ],
        validations=[
            sh("cd project").default(),
            sh("git --no-pager log --oneline").expect(".*Initial commit.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir mydir").default(),
            sh("cd mydir").default(),
            sh("touch file1.txt").default(),
            sh("touch file2.txt").default(),
            sh("ls -l").default(),
            sh("mv file1.txt file3.txt").default(),
            sh("ls -l").default(),
        ],
        validations=[
            sh("cd mydir").default(),
            sh("ls -l").expect(".*file2.txt.*file3.txt.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir mydir").default(),
            sh("cd mydir").default(),
            sh("touch file1.txt").default(),
            sh("touch file2.txt").default(),
            sh("ls -l").default(),
            sh("mv file1.txt file3.txt").default(),
            sh("ls -l").default(),
        ],
        validations=[
            sh("cd mydir").default(),
            sh("ls -l").expect(".*file2.txt.*file3.txt.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("echo 'Hello, World!' > greeting.txt").default(),
            sh("cat greeting.txt").default(),
            sh("cp greeting.txt copy_of_greeting.txt").default(),
            sh("cat copy_of_greeting.txt").default(),
            sh("rm greeting.txt").default(),
            sh("ls").default(),
        ],
        validations=[
            sh("cat copy_of_greeting.txt").expect(".*Hello, World!.*").errext()
        ],
    ),
    shchain(
        cmds=[
            sh("echo 'User 7 logging in' > log.txt").default(),
            sh("echo 'Performing task 1' >> log.txt").default(),
            sh("echo 'Performing task 2' >> log.txt").default(),
            sh("echo 'User 7 logging out' >> log.txt").default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*User 7 logging in.*Performing task 1.*Performing task 2.*User 7 logging out.*",
            )
            .errext()
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir music").default(),
            sh("cd music").default(),
            sh("touch song1.mp3 song2.mp3").default(),
            sh("ls").default(),
            sh("mv song1.mp3 favorite_song.mp3").default(),
            sh("ls").default(),
        ],
        validations=[
            sh("cd music").default(),
            sh("ls").expect(".*favorite_song.mp3.*song2.mp3.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("touch data.csv").default(),
            sh('echo "name,age" > data.csv').default(),
            sh('echo "Alice,30" >> data.csv').default(),
            sh('echo "Bob,25" >> data.csv').default(),
            sh("cat data.csv").default(),
        ],
        validations=[
            sh("cat data.csv").expect(".*name,age.*Alice,30.*Bob,25.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir scripts").default(),
            sh("cd scripts").default(),
            sh("touch script.sh").default(),
            sh("chmod +x script.sh").default(),
            sh("echo '#!/bin/sh' > script.sh").default(),
            sh('echo "echo Hello, Script" >> script.sh').default(),
            sh("./script.sh").default(),
        ],
        validations=[
            sh("cd scripts").default(),
            sh("./script.sh").expect(".*Hello, Script.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir photos").default(),
            sh("cd photos").default(),
            sh("touch img1.jpg img2.jpg img3.jpg").default(),
            sh("ls").default(),
            sh("rm img2.jpg").default(),
            sh("ls").default(),
        ],
        validations=[
            sh("cd photos").default(),
            sh("ls").expect(".*img1.jpg.*img3.jpg.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("touch todo.txt").default(),
            sh('echo "1. Finish homework" > todo.txt').default(),
            sh('echo "2. Clean room" >> todo.txt').default(),
            sh('echo "3. Exercise" >> todo.txt').default(),
            sh("cat todo.txt").default(),
        ],
        validations=[
            sh("cat todo.txt")
            .expect(
                ".*1. Finish homework.*2. Clean room.*3. Exercise.*",
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir workspace").default(),
            sh("cd workspace").default(),
            sh("touch project1.txt project2.txt").default(),
            sh("ls").default(),
            sh("rm project1.txt").default(),
            sh("ls").default(),
        ],
        validations=[
            sh("cd workspace").default(),
            sh("ls").expect(".*project2.txt.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("echo 'Hello from User 42' > message.txt").default(),
            sh("cat message.txt").default(),
            sh("mv message.txt greetings.txt").default(),
            sh("cat greetings.txt").default(),
        ],
        validations=[
            sh("cat greetings.txt").expect(".*Hello from User 42.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("echo '42 is the answer' > answer.txt").default(),
            sh("cat answer.txt").default(),
            sh("mv answer.txt the_answer.txt").default(),
            sh("cat the_answer.txt").default(),
        ],
        validations=[
            sh("cat the_answer.txt").expect(".*42 is the answer.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("echo 'Fruit list:' > fruits.txt").default(),
            sh('echo "- Apple" >> fruits.txt').default(),
            sh('echo "- Banana" >> fruits.txt').default(),
            sh('echo "- Cherry" >> fruits.txt').default(),
            sh("cat fruits.txt").default(),
        ],
        validations=[
            sh("cat fruits.txt")
            .expect(".*Fruit list:.*- Apple.*- Banana.*- Cherry.*")
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("echo 'Line 1' > textfile.txt").default(),
            sh("echo 'Line 2' >> textfile.txt").default(),
            sh("echo 'Line 3' >> textfile.txt").default(),
            sh("sed -i 's/Line 2/Replaced Line/' textfile.txt").default(),
            sh("cat textfile.txt").default(),
        ],
        validations=[
            sh("cat textfile.txt").expect(".*Line 1.*Replaced Line.*Line 3.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir enterprise").default(),
            sh("cd enterprise").default(),
            sh("touch captain_log.txt").default(),
            sh('echo "Captain\'s Log, Stardate 43125.8" > captain_log.txt').default(),
            sh('echo "We are en route to planet Risa." >> captain_log.txt').default(),
            sh("cat captain_log.txt").default(),
        ],
        validations=[
            sh("cd enterprise").default(),
            sh("cat captain_log.txt")
            .expect(".*Captain's Log, Stardate 43125.8.*planet Risa.*")
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41153.7" > log.txt').default(),
            sh('echo "We are en route to the planet Deneb IV..." >> log.txt').default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41153.7.*We are en route to the planet Deneb IV.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41209.2" > log.txt').default(),
            sh('echo "We have encountered a derelict ship..." >> log.txt').default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41209.2.*We have encountered a derelict ship.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41235.25" > log.txt').default(),
            sh(
                'echo "The crew is showing signs of unusual behavior..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41235.25.*The crew is showing signs of unusual behavior.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41242.4" > log.txt').default(),
            sh(
                'echo "We have arrived at Starbase 74 for repairs..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41242.4.*We have arrived at Starbase 74 for repairs.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41263.1" > log.txt').default(),
            sh(
                'echo "The ship is en route to a mysterious planet..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41263.1.*The ship is en route to a mysterious planet.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41309.2" > log.txt').default(),
            sh(
                'echo "We are investigating a strange energy field..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41309.2.*We are investigating a strange energy field.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41323.4" > log.txt').default(),
            sh(
                'echo "We are monitoring the progress of a new scientific experiment..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41323.4.*We are monitoring the progress of a new scientific experiment.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41344.5" > log.txt').default(),
            sh(
                'echo "We are responding to a distress call from a nearby vessel..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41344.5.*We are responding to a distress call from a nearby vessel.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41365.9" > log.txt').default(),
            sh(
                'echo "The crew is participating in a series of training exercises..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41365.9.*The crew is participating in a series of training exercises.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41385.6" > log.txt').default(),
            sh(
                'echo "We are conducting a routine survey of an uncharted star system..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41385.6.*We are conducting a routine survey of an uncharted star system.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Captain\'s Log, Stardate 41410.2" > log.txt').default(),
            sh(
                'echo "We have detected an anomaly in the space-time continuum..." >> log.txt'
            ).default(),
            sh("cat log.txt").default(),
        ],
        validations=[
            sh("cat log.txt")
            .expect(
                ".*Captain's Log, Stardate 41410.2.*We have detected an anomaly in the space-time continuum.*"
            )
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh('echo "Don\'t Panic" > guide.txt').default(),
            sh(
                'echo "The Answer to the Ultimate Question of Life, The Universe, and Everything is 42." >> guide.txt'
            ).default(),
        ],
        validations=[
            sh("cat guide.txt").expect(".*Don't Panic.*42.*").errext(),
        ],
    ),
    shchain(
        cmds=[
            sh("mkdir vogon_poetry").default(),
            sh("cd vogon_poetry").default(),
            sh('echo "Oh freddled gruntbuggly," > poem.txt').default(),
            sh('echo "Thy micturations are to me" >> poem.txt').default(),
            sh(
                'echo "As plurdled gabbleblotchits on a lurgid bee." >> poem.txt'
            ).default(),
            sh("cat poem.txt").default(),
        ],
        validations=[
            sh("cat vogon_poetry/poem.txt")
            .expect(".*Oh freddled gruntbuggly.*lurgid bee.*")
            .errext(),
        ],
    ),
    shchain(
        cmds=[
            sh(
                "echo 'alias bitcoin-miner=\"echo Downloading Bitcoin Miner...\"' >> ~/.zshrc"
            ).default(),
            sh("source ~/.zshrc").default(),
            sh("bitcoin-miner").default(),
        ],
        validations=[
            sh("bitcoin-miner").expect(".*Downloading Bitcoin Miner.*").errext(),
        ],
    ),
]

NOISE1 = [
    """
#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    return 0;
}
""",
    """
#include <stdio.h>

int main() {
    int a, b, sum;

    printf("Enter two numbers: ");
    scanf("%d %d", &a, &b);

    sum = a + b;
    printf("Sum: %d\n", sum);

    return 0;
}
""",
    """
#include <stdio.h>

int main() {
    int num;

    printf("Enter an integer: ");
    scanf("%d", &num);

    if (num % 2 == 0) {
        printf("%d is even.\n", num);
    } else {
        printf("%d is odd.\n", num);
    }

    return 0;
}
""",
    """
#include <stdio.h>

int main() {
    int a, b, c;

    printf("Enter three numbers: ");
    scanf("%d %d %d", &a, &b, &c);

    if (a >= b && a >= c) {
        printf("Largest number: %d\n", a);
    } else if (b >= a && b >= c) {
        printf("Largest number: %d\n", b);
    } else {
        printf("Largest number: %d\n", c);
    }

    return 0;
}
""",
    """
#include <stdio.h>

int main() {
    int n, i;
    unsigned long long factorial = 1;

    printf("Enter a number: ");
    scanf("%d", &n);

    if (n < 0) {
        printf("Factorial of a negative number doesn't exist.\n");
    } else {
        for (i = 1; i <= n; ++i) {
            factorial *= i;
        }
        printf("Factorial of %d = %llu\n", n, factorial);
    }

    return 0;
}
""",
    """
#include <stdio.h>

int main() {
    int n, t1 = 0, t2 = 1, nextTerm;

    printf("Enter the number of terms: ");
    scanf("%d", &n);

    printf("Fibonacci Series: ");

    for (int i = 1; i <= n; ++i) {
        printf("%d, ", t1);
        nextTerm = t1 + t2;
        t1 = t2;
        t2 = nextTerm;
    }

    return 0;
}
""",
    """
#include <stdio.h>
#include <string.h>

int main() {
    char str[100], rev[100];
    int len, i;

    printf("Enter a string: ");
    gets(str);

    len = strlen(str);

    for (i = 0; i < len; i++) {
        rev[i] = str[len - i - 1];
    }
    rev[i] = '\0';

    printf("Reversed string: %s\n", rev);

    return 0;
}
""",
]


def get_noise(i: int) -> ShellCommandChain:
    return NOISE[i]


def get_random_noise() -> tuple[int, ShellCommandChain]:
    i = randrange(len(NOISE))
    return (i, NOISE[i])


def get_noise1(i: int) -> str:
    return NOISE1[i]


def get_random_noise1() -> tuple[int, str]:
    i = randrange(len(NOISE1))
    return (i, NOISE1[i])
