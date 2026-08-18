#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#include "mongoose.h"
#include "db.h"
#include "http.h"

volatile sig_atomic_t run = 1;

void signal_handler(int sig) {
	(void) sig;
	run = 0;
}

int main(int argc, char** argv) {
	struct mg_mgr mgr;

	signal(SIGINT, signal_handler);
	signal(SIGTERM, signal_handler);


	if (argc < 2) {
		fprintf(stderr, "usage: %s <ADDR>\n", argv[1]);
		return 1;
	}

	if (db_init()) {
		fprintf(stderr, "failed to initialize db\n");
		return 1;
	}

	mg_mgr_init(&mgr);

	if (!mg_http_listen(&mgr, argv[1], ev_handler, NULL)) {
		fprintf(stderr, "failed to listen to %s\n", argv[1]);
		db_fini();
		mg_mgr_free(&mgr);
		return 1;
	}

	while (run)
		mg_mgr_poll(&mgr, 10);

	db_fini();
	mg_mgr_free(&mgr);
}
