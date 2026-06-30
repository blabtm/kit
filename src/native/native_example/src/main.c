#include <stdbool.h>
#include <stdio.h>
#include <unistd.h>

int main(void) {
  while (true) {
    printf("I'm alive!\n");
    sleep(3);
  }

  return 0;
}
