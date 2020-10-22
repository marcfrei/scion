/*
 * go build -o bin/libhello_c.a -buildmode=c-archive go/experiments/hello/hello_c.go
 * cc -o bin/hello_c -Wall -Werror -I bin -L bin go/experiments/hello/hello_c.c -lhello_c -lpthread
 */

#include "libhello_c.h"

int
main(int argc, char * const *argv)
{
	if (argc < 4) {
		return 1;
	}

	(void)SendHello(argv[1], argv[2], argv[3]);
	return 0;
}
