#include <signal.h>

#include "server.h"
#include "embed.h"

void r_index(struct request* req, struct response* res) {
	res_write(res, INDEX, sizeof(INDEX));
	res_status(res, STATUS_OK);
}

void r_page2(struct request* req, struct response* res) {
	res_write(res, PAGE2, sizeof(PAGE2));
	res_status(res, STATUS_OK);
}

void r_page3(struct request* req, struct response* res) {
	res_write(res, PAGE3, sizeof(PAGE3));
	res_status(res, STATUS_OK);
}

int main(int argc, char** argv) {
	signal(SIGPIPE, SIG_IGN);
	server_add("/", &r_index);
	server_add("/index.html", &r_index);
	server_add("/page2.html", &r_page2);
	server_add("/page3.html", &r_page3);
	server_listen(8082);
}
