#include "db.h"
#include "web.h"
#include <signal.h>
#include <stdio.h>

static int sig_exit = 0;
static void signal_handler(int sig) {
	printf("signal %d received\n", sig);
	sig_exit = 1;
}

int main(void) {
	db_init();
	web_init();

	signal(SIGINT, signal_handler);
	signal(SIGTERM, signal_handler);
	signal(SIGKILL, signal_handler);
	signal(SIGQUIT, signal_handler);
	signal(SIGABRT, signal_handler);

	while(sig_exit == 0)
		web_step();

	web_fini();
	db_fini();
}
