#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <fcntl.h>

#include "server.h"

struct header {
	char key[64];
	char val[1024];
};

struct request {
	int  fd;
	char method[8];
	char uri[1024];
	char version[16];
	struct header headers[32];
	int header_count;
};

struct response {
	int  fd;
	int  status;
	char buf[8192];
	int  len;
};

struct route {
	const char* path;
	void (*f)(struct request*, struct response*);
};

static struct route* routes = NULL;
int routes_len = 0, routes_cap = 0;

void server_add(const char* path, void (*f)(struct request*, struct response*)) {
	if (routes_len + 1 > routes_cap) {
		int new_cap = routes_cap == 0 ? 4 : 2 * routes_cap;
		routes = realloc(routes, new_cap * sizeof(struct route));
		routes_cap = new_cap;
	}

	struct route* route = &routes[routes_len++];
	route->path = path;
	route->f = f;
}

static void server_dispatch(struct request* req, struct response* res) {
	for (size_t i = 0; i < routes_len; i++)
		if (strcmp(req->uri, routes[i].path) == 0) {
			routes[i].f(req, res);
			return;
		}
	res_status(res, STATUS_NOT_FOUND);
}

static int req_parse(struct request* req, int fd) {
	char buf[8192];
	FILE* fp = fdopen(fd, "r");
	char* r = fgets(buf, sizeof(buf), fp);
	if (r == NULL)
		return -1;
	if (sscanf(buf, "%7s %1023s %15s", req->method, req->uri, req->version) != 3)
		return -1;

	while (1) {
		r = fgets(buf, sizeof(buf), fp);
		if (r == NULL)
			return -1;
		if (buf[0] == '\r' && buf[1] == '\n')
			break;
		struct header* header = &req->headers[req->header_count++];
		if (sscanf(buf, "%63[^:]: %1023[^\r\n]", header->key, header->val) != 2)
			return -1;
	}

	return 0;
}

void server_listen(int port) {
	int sfd = socket(AF_INET, SOCK_STREAM, 0);
	if (sfd == -1)
		perror("socket"), exit(1);

	struct sockaddr_in addr;
	int addr_len = sizeof(addr);
	memset(&addr, 0, addr_len);
	addr.sin_family      = AF_INET;
	addr.sin_port        = htons(port);
	addr.sin_addr.s_addr = htonl(INADDR_ANY);

	if (bind(sfd, (struct sockaddr*) &addr, addr_len) == -1)
		perror("bind"), exit(1);

	if (listen(sfd, SOMAXCONN) == -1)
		perror("listen"), exit(1);

	printf("server listening on :%d\n", port);
	while (1) {
		int rfd = accept(sfd, (struct sockaddr*) &addr, (socklen_t*) &addr_len);
		if (rfd == -1) {
			perror("accept");
			continue;
		}
		printf("received request\n");

		struct request  req = { 0 };
		struct response res = { 0 };
		res.fd = rfd;
		if (req_parse(&req, rfd) != 0) {
			res_status(&res, STATUS_BAD_REQUEST);
			close(rfd);
			continue;
		}
		printf("parsed request\n");
		server_dispatch(&req, &res);
		printf("did dispatch\n");

		int r;
		if (res.status == STATUS_OK)
			r = write(res.fd, "HTTP/1.1 200 OK\r\n", strlen("HTTP/1.1 200 OK\r\n"));
		if (res.status == STATUS_NOT_FOUND)
			r = write(res.fd, "HTTP/1.1 404 NOT FOUND\r\n", strlen("HTTP/1.1 404 NOT FOUND\r\n"));
		if (r == -1) {
			close(rfd);
			continue;
		}
		char buf[8192];
		snprintf(buf, sizeof(buf), "Content-Length: %d\r\n", res.len);
		r = write(res.fd, buf, strlen(buf));
		if (r == -1) {
			close(rfd);
			continue;
		}
		r = write(res.fd, "\r\n", 2);
		if (r == -1) {
			close(rfd);
			continue;
		}
		r = write(res.fd, res.buf, res.len);
		if (r == -1) {
			close(rfd);
			continue;
		}

		close(rfd);
	}
	close(sfd);
}

void res_write(struct response* res, const char* buf, int len) {
	memcpy(&res->buf[res->len], buf, len);
	res->len += len;
}

void res_status(struct response* res, http_status status) {
	res->status = status;
}
