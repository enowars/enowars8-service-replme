from random import randrange

from util import ShellCommandChain, shchain, sh

NOISE: list[ShellCommandChain] = [
    shchain(
        cmds=[
            sh('echo "2b or not 2b" > note.txt').default(),
        ],
        validations=[
            sh("cat note.txt").expect(".*\n2b or not 2b.*").errext(),
        ],
    ),
    # shchain(
    #     cmds=[
    #         sh('echo "#include <stdio.h>" > main.c').default(),
    #         sh('echo "int main() {" >> main.c').default(),
    #         sh("echo '    printf(\"Hello World!\\\\n\");' >> main.c").default(),
    #         sh('echo "    return 0;" >> main.c').default(),
    #         sh('echo "}" >> main.c').default(),
    #         sh("gcc -o main main.c").default(),
    #         sh("./main").expect(".*\nHello World!.*").default(),
    #     ],
    #     validations=[
    #         sh("gcc -o main main.c").default(),
    #         sh("./main").expect(".*\nHello World!.*").default(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir docs").default(),
    #         sh("cd docs").default(),
    #         sh("ls -lah").default(),
    #         sh('echo "And I say hey"').default(),
    #         sh('echo "What\'s going on"').default(),
    #         sh("touch 4-non-blondes.txt").default(),
    #         sh('echo "HEY, YEAH, YEAH" >> 4-non-blondes.txt').default(),
    #     ],
    #     validations=[
    #         sh("cd docs").default(),
    #         sh("cat 4-non-blondes.txt").expect(".*\nHEY, YEAH, YEAH.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("touch notes.txt").default(),
    #         sh('echo "To do list:" > notes.txt').default(),
    #         sh('echo "- Buy groceries" >> notes.txt').default(),
    #         sh('echo "- Call mom" >> notes.txt').default(),
    #         sh('echo "- Finish project" >> notes.txt').default(),
    #         sh("cat notes.txt").default(),
    #     ],
    #     validations=[
    #         sh("cat notes.txt")
    #         .expect(
    #             ".*\nTo do list:.*\n- Buy groceries.*\n- Call mom.*\n- Finish project.*",
    #         )
    #         .errext()
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir project").default(),
    #         sh("cd project").default(),
    #         sh('git config --global user.email "1337@pwn.y"').default(),
    #         sh('git config --global user.name "Eve"').default(),
    #         sh("git init").default(),
    #         sh('echo "# My Project" > README.md').default(),
    #         sh("git add README.md").default(),
    #         sh('git commit -m "Initial commit"').default(),
    #     ],
    #     validations=[
    #         sh("cd project").default(),
    #         sh("git --no-pager log --oneline").expect(".*Initial commit.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir mydir").default(),
    #         sh("cd mydir").default(),
    #         sh("touch file1.txt").default(),
    #         sh("touch file2.txt").default(),
    #         sh("ls -l").default(),
    #         sh("mv file1.txt file3.txt").default(),
    #         sh("ls -l").default(),
    #     ],
    #     validations=[
    #         sh("cd mydir").default(),
    #         sh("ls -l").expect(".*file2.txt.*file3.txt.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir mydir").default(),
    #         sh("cd mydir").default(),
    #         sh("touch file1.txt").default(),
    #         sh("touch file2.txt").default(),
    #         sh("ls -l").default(),
    #         sh("mv file1.txt file3.txt").default(),
    #         sh("ls -l").default(),
    #     ],
    #     validations=[
    #         sh("cd mydir").default(),
    #         sh("ls -l").expect(".*file2.txt.*file3.txt.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("echo 'Hello, World!' > greeting.txt").default(),
    #         sh("cat greeting.txt").default(),
    #         sh("cp greeting.txt copy_of_greeting.txt").default(),
    #         sh("cat copy_of_greeting.txt").default(),
    #         sh("rm greeting.txt").default(),
    #         sh("ls").default(),
    #     ],
    #     validations=[
    #         sh("cat copy_of_greeting.txt").expect(".*Hello, World!.*").errext()
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("echo 'User 7 logging in' > log.txt").default(),
    #         sh("echo 'Performing task 1' >> log.txt").default(),
    #         sh("echo 'Performing task 2' >> log.txt").default(),
    #         sh("echo 'User 7 logging out' >> log.txt").default(),
    #         sh("cat log.txt").default(),
    #     ],
    #     validations=[
    #         sh("cat log.txt")
    #         .expect(
    #             ".*User 7 logging in.*Performing task 1.*Performing task 2.*User 7 logging out.*",
    #         )
    #         .errext()
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir music").default(),
    #         sh("cd music").default(),
    #         sh("touch song1.mp3 song2.mp3").default(),
    #         sh("ls").default(),
    #         sh("mv song1.mp3 favorite_song.mp3").default(),
    #         sh("ls").default(),
    #     ],
    #     validations=[
    #         sh("cd music").default(),
    #         sh("ls").expect(".*favorite_song.mp3.*song2.mp3.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("touch data.csv").default(),
    #         sh('echo "name,age" > data.csv').default(),
    #         sh('echo "Alice,30" >> data.csv').default(),
    #         sh('echo "Bob,25" >> data.csv').default(),
    #         sh("cat data.csv").default(),
    #     ],
    #     validations=[
    #         sh("cat data.csv").expect(".*name,age.*Alice,30.*Bob,25.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir scripts").default(),
    #         sh("cd scripts").default(),
    #         sh("touch script.sh").default(),
    #         sh("chmod +x script.sh").default(),
    #         sh("echo '#!/bin/sh' > script.sh").default(),
    #         sh('echo "echo Hello, Script" >> script.sh').default(),
    #         sh("./script.sh").default(),
    #     ],
    #     validations=[
    #         sh("cd scripts").default(),
    #         sh("./script.sh").expect(".*Hello, Script.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir photos").default(),
    #         sh("cd photos").default(),
    #         sh("touch img1.jpg img2.jpg img3.jpg").default(),
    #         sh("ls").default(),
    #         sh("rm img2.jpg").default(),
    #         sh("ls").default(),
    #     ],
    #     validations=[
    #         sh("cd photos").default(),
    #         sh("ls").expect(".*img1.jpg.*img3.jpg.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("touch todo.txt").default(),
    #         sh('echo "1. Finish homework" > todo.txt').default(),
    #         sh('echo "2. Clean room" >> todo.txt').default(),
    #         sh('echo "3. Exercise" >> todo.txt').default(),
    #         sh("cat todo.txt").default(),
    #     ],
    #     validations=[
    #         sh("cat todo.txt")
    #         .expect(
    #             ".*1. Finish homework.*2. Clean room.*3. Exercise.*",
    #         )
    #         .errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("mkdir workspace").default(),
    #         sh("cd workspace").default(),
    #         sh("touch project1.txt project2.txt").default(),
    #         sh("ls").default(),
    #         sh("rm project1.txt").default(),
    #         sh("ls").default(),
    #     ],
    #     validations=[
    #         sh("cd workspace").default(),
    #         sh("ls").expect(".*project2.txt.*").errext(),
    #     ],
    # ),
    # shchain(
    #     cmds=[
    #         sh("echo 'Hello from User 42' > message.txt").default(),
    #         sh("cat message.txt").default(),
    #         sh("mv message.txt greetings.txt").default(),
    #         sh("cat greetings.txt").default(),
    #     ],
    #     validations=[
    #         sh("cat greetings.txt").expect(".*Hello from User 42.*").errext(),
    #     ],
    # ),
]


def get_noise(i: int) -> ShellCommandChain:
    return NOISE[i]


def get_random_noise() -> tuple[int, ShellCommandChain]:
    i = randrange(len(NOISE))
    return (i, NOISE[i])
