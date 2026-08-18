#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#include "mongoose.h"
#include "db.h"
#include "http.h"

static volatile sig_atomic_t run = 1;

static void signal_handler(int sig) {
	(void) sig;
	run = 0;
}

int main(int argc, char** argv) {
	struct mg_mgr mgr;

	signal(SIGINT, signal_handler);
	signal(SIGTERM, signal_handler);

	db_init();
	mg_mgr_init(&mgr);
	mg_http_listen(&mgr, argv[1], ev_handler, NULL);

	while (run)
		mg_mgr_poll(&mgr, 10);

	db_fini();
	mg_mgr_free(&mgr);
}
