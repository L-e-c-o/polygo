// Modified version of https://gist.github.com/securitytube/5318838#gistcomment-3614284
// gcc -m32 -Wall runnner.c -o runner

#include <sys/mman.h>

unsigned char code[] = "\x90\x90\x90\x90";

int main(){
	mprotect((void *)((int)code & ~4095),  4096, PROT_READ | PROT_WRITE | PROT_EXEC);
	int (*ret)() = (int(*)())code;
	return ret();
}