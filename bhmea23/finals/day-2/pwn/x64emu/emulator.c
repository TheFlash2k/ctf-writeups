#include <assert.h>
#include <setjmp.h>
#include <signal.h>
#include <stdio.h>
#include <sys/time.h>
#include <unicorn/unicorn.h>

#define CODE_BASE 0x400000UL
#define CODE_SIZE 0x010000UL
#define TIMEOUT_MAX (3*1000*1000)

jmp_buf g_env;

void timeout_handler(int sig) {
  puts("Timeout.");
  siglongjmp(g_env, 1);
}

void setup_handler() {
  struct sigaction sa;
  sa.sa_handler = timeout_handler;
  sigemptyset(&(sa.sa_mask));
  sigaddset(&(sa.sa_mask), SIGALRM);
  sigaction(SIGALRM, &sa, NULL);
}

int main() {
  struct itimerval timeout;
  uc_engine *uc;
  uc_hook tr;
  size_t size, usec, usec_remain;
  uint8_t *code;

  setvbuf(stdin, NULL, _IONBF, 0);
  setvbuf(stdout, NULL, _IONBF, 0);
  setvbuf(stderr, NULL, _IONBF, 0);

  /* Setup unicorn */
  if (sigsetjmp(g_env, 1) == 0) {
    assert (uc_open(UC_ARCH_X86, UC_MODE_64, &uc) == 0);
    uc_mem_map(uc, CODE_BASE, CODE_SIZE, UC_PROT_READ | UC_PROT_EXEC);
  }

  /* Input timeout */
  printf("timeout (usec): ");
  assert (scanf("%ld%*c", &usec) == 1);
  assert (0 < usec && usec <= TIMEOUT_MAX);
  timeout.it_interval.tv_sec = 0;
  timeout.it_interval.tv_usec = 0;
  timeout.it_value.tv_sec = usec / (1000*1000);
  timeout.it_value.tv_usec = usec % (1000*1000);

  /* Input code */
  printf("size: ");
  assert (scanf("%ld%*c", &size) == 1);
  assert (size <= CODE_SIZE);
  printf("code: ");
  assert (code = (uint8_t*)malloc(size));
  assert (fread(code, sizeof(uint8_t), size, stdin) == size);
  assert (uc_mem_write(uc, CODE_BASE, code, size) == 0);
  free(code);

  setup_handler();

  /* Run emulator */
  assert (setitimer(ITIMER_REAL, &timeout, NULL) == 0);
  uc_emu_start(uc, CODE_BASE, CODE_BASE + CODE_SIZE, 0, 0);
  uc_close(uc);
  puts("Emulation done.");

  /* Calculate total time */
  assert (getitimer(ITIMER_REAL, &timeout) == 0);
  usec_remain = timeout.it_value.tv_sec*(1000*1000) + timeout.it_value.tv_usec;
  printf("Elapsed: %ld usec\n", usec - usec_remain);

  return 0;
}
