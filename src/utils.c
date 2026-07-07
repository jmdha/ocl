#include <stdlib.h>
#include <unistd.h>
#include <sys/stat.h>

#include "utils.h"

int ranfile(FILE** f, char* path, size_t size, const char* dir) {
	int  n;
	int  w;
	
	mkdir(dir, 0777);

	while (1) {
		n = rand();
		w = snprintf(path, size, "%s/%d", dir, n);
		if (w < 0) {
			perror("ranfile(name)");
			return -1;
		}
		if (access(path, F_OK) == 0)
			continue;
		*f = fopen(path, "w");
		if (*f == NULL) {
			perror("ranfile(fopen)");
			return -1;
		}
		break;
	}

	return 0;
}

