/*
 * go build -buildmode=c-archive -o bin/libhello_s.a go/experiments/hello/hello_s.go go/experiments/hello/hello_s_api.go
 * cc -o bin/hello_s -Wall -Werror -I bin -L bin go/experiments/hello/hello_s.c -lhello_s -lpthread
 */

#include <stdio.h>
#include "libhello_s.h"

void
handle_data(char* data)
{
	printf("Received data: \"%s\"\n", data);
}

int
main(int argc, char * const *argv)
{
	if (argc < 3) {
		return 1;
	}

	(void)RunServer(argv[1], argv[2], handle_data);
	return 0;
}
