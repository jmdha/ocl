#include <signal.h>
#include <stdbool.h>

#include "db.h"
#include "web.h"

static volatile sig_atomic_t running = 1;

static void signal_handler(int sig) {
    (void) sig;
    running = 0;
}

int main(int argc, char** argv) {
	signal(SIGINT, signal_handler);   // Ctrl+C
	signal(SIGTERM, signal_handler);  // kill <pid>
					  
	if (db_init() != 0) {
		fprintf(stderr, "failed to init db\n");
		return 1;
	}

	if (web_init(atoi(argv[1])) != 0) {
		fprintf(stderr, "failed to init web\n");
		return 1;
	}

	while (running)
		web_poll();
	db_fini();
	web_fini();
}
