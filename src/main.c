#include <signal.h>
#include <stdio.h>
#include <stdlib.h>

#include "db.h"
#include "http.h"
#include "routes.h"

int main(int argc, char** argv) {
	// a peer closing mid response must not kill the server
	signal(SIGPIPE, SIG_IGN);

	if (argc < 2) {
		fprintf(stderr, "usage: %s <PORT>\n", argv[0]);
		return 1;
	}

	if (db_init()) {
		fprintf(stderr, "failed to initialize db\n");
		return 1;
	}

	http_register("POST", "/api/public/users", post_users);
	http_register("POST", "/api/public/login", post_login);
	http_register("GET",  "/api/upload",       get_upload);
	http_register("POST", "/api/upload",       post_upload);
	http_listen(atoi(argv[1]));
}
