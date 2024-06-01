#include <pthread.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*
CRC-32
poly       init       refin refout xorout
0x04C11DB7 0xFFFFFFFF true  true   0xFFFFFFFF

crc32("test")=0xd87f7e0c = reverse(0xCF8101E4 ^ 0xFFFFFFFF)

aacuwzj
botealu
m}eyvtl
nsri`bs
r|algxe
qrv|qnz
~`g`fvc
}nppp`|


CRC-32/BZIP2 (kinda)
poly       init       refin refout xorout
0x04C11DB7 0x00000000 false false  0xFFFFFFFF

crc32("test")=0xf48f12d7=0xb70ed28^0xFFFFFFFF

bau~gv`
cx`lksa
dfyxwxl
}qqkvpr
|hdyzus
{v}mf~~

*/
unsigned char reverse8(unsigned char b) {
  b = (b & 0xF0) >> 4 | (b & 0x0F) << 4;
  b = (b & 0xCC) >> 2 | (b & 0x33) << 2;
  b = (b & 0xAA) >> 1 | (b & 0x55) << 1;
  return b;
}

uint32_t reverse32(uint32_t x) {
  x = ((x >> 1) & 0x55555555u) | ((x & 0x55555555u) << 1);
  x = ((x >> 2) & 0x33333333u) | ((x & 0x33333333u) << 2);
  x = ((x >> 4) & 0x0f0f0f0fu) | ((x & 0x0f0f0f0fu) << 4);
  x = ((x >> 8) & 0x00ff00ffu) | ((x & 0x00ff00ffu) << 8);
  x = ((x >> 16) & 0xffffu) | ((x & 0xffffu) << 16);
  return x;
}

uint8_t crc8(char *buf, size_t len) {
  uint8_t poly = 0b11010101;
  uint8_t reg = 0x00;
  uint8_t xor = 0x00;
  for (int j = 0; j < len; j++) {
    unsigned char b = buf[j];
    for (int i = 0; i < 8; i++) {
      if ((reg ^ b) & 0x80) {
        reg = (reg << 1) ^ poly;
      } else {
        reg <<= 1;
      }
      b <<= 1;
    }
  }
  return reg ^ xor;
}

void brute_force_crc8() {
  uint8_t buf[3] = {0};
  size_t preimage_len = sizeof(buf) - 1;
  char *goal_str = "aa";
  uint8_t goal = crc8(goal_str, strlen(goal_str));
  printf("crc8(\"%s\")=0x%x\n", goal_str, goal);
  for (int i = 1; i <= preimage_len; i++) {
    for (int j = 0; j < i; j++)
      buf[j] = 0b01100000; // `
    buf[i] = '\0';
    while (1) {
      if (crc8((char *)buf, i) == goal)
        printf(" Match: %s\n", buf);
      int j = 0;
      for (; j < i; j++) {
        if (buf[j] == 0b01111110) { // ~
          buf[j] = 0b01100000;      // `
        } else {
          buf[j]++;
          break;
        }
      }
      if (j == i)
        break;
    }
  }
}

// CRC-32/BZIP2
// poly       init       refin refout xorout
// 0x04C11DB7 0xFFFFFFFF false false  0xFFFFFFFF
uint32_t crc32(char *buf, size_t len) {
  uint32_t poly = 0x04C11DB7;
  uint32_t xor = 0xFFFFFFFF;
  uint32_t reg = 0x0;
  for (int j = 0; j < len; j++) {
    unsigned char b = buf[j];
    for (int i = 0; i < 8; i++) {
      if ((reg ^ (b << 24)) & 0x80000000) {
        reg = (reg << 1) ^ poly;
      } else {
        reg <<= 1;
      }
      b <<= 1;
    }
  }
  return (reg) ^ xor;
}

// CRC-32
// poly       init       refin refout xorout
// 0x04C11DB7 0xFFFFFFFF true  true   0xFFFFFFFF
uint32_t crc32_rev(char *buf, size_t len) {
  uint32_t poly = 0x04C11DB7;
  uint32_t xor = 0xFFFFFFFF;
  uint32_t reg = 0xFFFFFFFF;
  for (int j = 0; j < len; j++) {
    unsigned char b = reverse8(buf[j]);
    for (int i = 0; i < 8; i++) {
      if ((reg ^ (b << 24)) & 0x80000000) {
        reg = (reg << 1) ^ poly;
      } else {
        reg <<= 1;
      }
      b <<= 1;
    }
  }
  printf("0x%x\n", reg);
  return reverse32(reg ^ xor);
}

uint32_t *crc32_rev_table(uint32_t poly) {
  uint32_t *table = malloc(sizeof(uint32_t) * 256);
  for (uint16_t i = 0; i <= 0xff; i++) {
    uint8_t b = reverse8(i & 0xff);
    uint32_t reg = b << 24;
    for (int i = 0; i < 8; i++) {
      if (reg & 0x80000000) {
        reg = (reg << 1) ^ poly;
      } else {
        reg <<= 1;
      }
    }
    table[i] = reverse32(reg);
  }
  return table;
}

uint32_t crc32_with_table(uint32_t *table, char *buf, size_t len) {
  uint32_t xor = 0xFFFFFFFF;
  uint32_t reg = reverse32(0xFFFFFFFF);
  for (int j = 0; j < len; j++) {
    reg = table[(reg ^ buf[j]) & 0xff] ^ (reg >> 8);
  }
  return reg ^ xor;
}

struct bf_crc32_thread_args {
  size_t iter;
  size_t init;
  size_t step;
  uint32_t goal;
  uint32_t *table;
};

void *bf_crc32_table_thread(void *_args) {
  struct bf_crc32_thread_args *args = _args;
  size_t len = args->iter;
  uint8_t *buf = malloc(sizeof(uint8_t) * (len + 1));
  for (int j = 0; j < len; j++)
    buf[j] = 0b01100000; // `
  buf[len] = '\0';
  buf[0] += args->init;
  for (int i = 0; i <= (32 - args->step); i += args->step) {
    while (1) {
      if (crc32_with_table(args->table, (char *)buf, len) == args->goal)
        printf(" Match: %s\n", buf);
      int j = 1;
      for (; j < len; j++) {
        if (buf[j] >= 0b01111110) { // ~
          buf[j] = 0b01100000;      // `
        } else {
          buf[j]++;
          break;
        }
      }
      if (j == len)
        break;
    }
    buf[0] += args->step;
  }

  free(buf);

  return NULL;
}

void *bf_crc32_thread(void *_args) {
  struct bf_crc32_thread_args *args = _args;
  size_t len = args->iter;
  uint8_t *buf = malloc(sizeof(uint8_t) * (len + 1));
  for (int j = 0; j < len; j++)
    buf[j] = 0b01100000; // `
  buf[len] = '\0';
  buf[0] += args->init;
  for (int i = 0; i <= (32 - args->step); i += args->step) {
    while (1) {
      if (crc32((char *)buf, len) == args->goal)
        printf(" Match: %s\n", buf);
      int j = 1;
      for (; j < len; j++) {
        if (buf[j] >= 0b01111110) { // ~
          buf[j] = 0b01100000;      // `
        } else {
          buf[j]++;
          break;
        }
      }
      if (j == len)
        break;
    }
    buf[0] += args->step;
  }

  free(buf);

  return NULL;
}

#define NUM_THREADS 8

void brute_force_crc32_table() {
  uint32_t *table = crc32_rev_table(0x04C11DB7);
  char *goal_str = "test";
  uint32_t goal = crc32_with_table(table, goal_str, strlen(goal_str));
  printf("crc32(\"%s\")=0x%x\n", goal_str, goal);
  printf("              0x%x\n", reverse32(goal ^ 0xffffffff));

  pthread_t ids[NUM_THREADS] = {0};
  struct bf_crc32_thread_args args[NUM_THREADS] = {0};

  for (int i = 0; i < NUM_THREADS; i++) {
    args[i].table = table;
  }

  for (int i = 1; i <= 8; i++) {
    printf("Len %d\n", i);
    for (int j = 0; j < NUM_THREADS; j++) {
      args[j].step = NUM_THREADS;
      args[j].init = j;
      args[j].iter = i;
      args[j].goal = goal;
      pthread_create(ids + j, NULL, bf_crc32_table_thread, args + j);
    }
    for (int j = 0; j < NUM_THREADS; j++) {
      pthread_join(ids[j], NULL);
    }
  }

  free(table);
}

void brute_force_crc32() {
  char *goal_str = "test";
  uint32_t goal = crc32(goal_str, strlen(goal_str));
  printf("crc32(\"%s\")=0x%x\n", goal_str, goal);
  printf("              0x%x\n", goal ^ 0xFFFFFFFF);

  pthread_t ids[NUM_THREADS] = {0};
  struct bf_crc32_thread_args args[NUM_THREADS] = {0};

  for (int i = 1; i <= 7; i++) {
    printf("Len %d\n", i);
    for (int j = 0; j < NUM_THREADS; j++) {
      args[j].step = NUM_THREADS;
      args[j].init = j;
      args[j].iter = i;
      args[j].goal = goal;
      pthread_create(ids + j, NULL, bf_crc32_thread, args + j);
    }
    for (int j = 0; j < NUM_THREADS; j++) {
      pthread_join(ids[j], NULL);
    }
  }
}

int main(int args, char *argv[]) {
  //
  // brute_force_crc32();
  printf("0x%x\n", crc8("aa", 2));
}
