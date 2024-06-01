#include <zlib.h>
#include <stdio.h>

int main() {
    int max_len = 16;
    // 65 - 90
    // 97 - 122

    char mask[26*2];

    for (int i = 0; i < 26; i++) {
        mask[i] = i + 'A';
        mask[i+26] = i + 'a';
    }

    char *s = "test";
    unsigned long goal = crc32(0, s, strlen(s));

    for (int i = 1; i <= 16; i++) {
        char buf[i];
        for (int j = 0; j < i; j++) {
            buf[j] = 'a';
        }
        while (1) {
            int j;
            for (j = 0; j < i; j++) {
                if (buf[j] < 'z') {
                    buf[j]++;
                    break;
                }
                if (buf[j] == 'z') {
                    buf[j] = 'a';
                }
            }
            if (j == i) {
                break;
            }

            unsigned long test = crc32(0, buf, i);
            if (test == goal)
                printf("Match: %.*s\n", i, buf);
        }
    }

    return 0;
}
