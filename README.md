# repl.me

Throwaway devenvs in the web.

## idea

Server manages docker containers. Streams io over ws connection between client and pty (bash) wrapper service running on each container.

Clients can access either shared container through user+pass, or let server spin up personal container identified through hash of username.

Clients on shared container can interact with each other in "coffeelounge". Simple program to allow some chatting.

## vulns

1. Hash function for identification of personal container is weak. Has multiple master images.
2. Shared container has misconfigurations allowing access of files of other users (TBD)
3. Some vuln in coffeelounge (TBD)

## flags

Flags are stored in some user's home on the shared and personal containers. Flags are not stored in cleartext, but are printed in stdout by running some brainfuck code. Brainfuck program needs some input to generate the correct output from. Input could be "ENO".

## stack

- pty wrapper: xterm.js for prototype; C, libvterm for mvp
- server: go / next.js (TBD)
- client: html/css/js / react.js / next.js (TBD)
- coffeelounge: java

