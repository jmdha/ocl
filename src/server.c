#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <fcntl.h>

#include "server.h"
#include "log.h"

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
	size_t bread;
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

int SOCKET = 0;
static struct route* routes = NULL;
int routes_len = 0, routes_cap = 0;

void fini() {
	if (SOCKET)
		close(SOCKET);
}

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
	int o = 0;
	while (1) {
		char c;
		int r = read(fd, &c, 1);
		if (r <= 0)
			return -1;
		if (o >= 8192)
			break;
		buf[o++] = c;
		if (c == '\n')
			break;
	}
	buf[o] = '\0';
	if (sscanf(buf, "%7s %1023s %15s", req->method, req->uri, req->version) != 3)
		return -1;

	while (1) {
		o = 0;
		memset(buf, 0, sizeof(buf));
		while (1) {
			char c;
			int r = read(fd, &c, 1);
			if (r <= 0)
				return -1;
			if (o >= 8192)
				break;
			buf[o++] = c;
			if (c == '\n')
				break;
		}
		buf[o] = '\0';
		if (buf[0] == '\r' && buf[1] == '\n')
			break;
		struct header* header = &req->headers[req->header_count++];
		if (sscanf(buf, "%63[^:]: %1023[^\r\n]", header->key, header->val) != 2) {
			printf("Invalid header\n");
			return -1;
		}
	}

	return 0;
}

void server_listen(int port) {
	SOCKET = socket(AF_INET, SOCK_STREAM, 0);
	if (SOCKET == -1)
		perror("socket"), exit(1);
	atexit(fini);

	struct sockaddr_in addr;
	int addr_len = sizeof(addr);
	memset(&addr, 0, addr_len);
	addr.sin_family      = AF_INET;
	addr.sin_port        = htons(port);
	addr.sin_addr.s_addr = htonl(INADDR_ANY);

	if (bind(SOCKET, (struct sockaddr*) &addr, addr_len) == -1)
		perror("bind"), exit(1);

	if (listen(SOCKET, SOMAXCONN) == -1)
		perror("listen"), exit(1);

	printf("server listening on :%d\n", port);
	while (1) {
		int rfd = accept(SOCKET, (struct sockaddr*) &addr, (socklen_t*) &addr_len);
		if (rfd == -1) {
			perror("accept");
			continue;
		}

		struct request  req = { 0 };
		struct response res = { 0 };
		res.fd = rfd;
		req.fd = rfd;
		if (req_parse(&req, rfd) != 0) {
			res_status(&res, STATUS_BAD_REQUEST);
			close(rfd);
			continue;
		}
		log_info("%s %s\n", req.method, req.uri);
		server_dispatch(&req, &res);

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
		snprintf(buf, sizeof(buf), "Connection: close\r\n");
		r = write(res.fd, buf, strlen(buf));
		if (r == -1) {
			close(rfd);
			continue;
		}
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
}

int req_read(struct request* req, char* buf, int cap) {
	for (size_t i = 0; i < req->header_count; i++)
		if (strcmp(req->headers[i].key, "Content-Length") == 0) {
			int tcap = atoi(req->headers[i].val);
			if (req->bread >= tcap)
				return 0;
			if (tcap - req->bread < cap)
				cap = tcap - req->bread;
			break;
		}
	int bread = read(req->fd, buf, cap);
	if (bread == -1)
		return -1;
	req->bread += bread;
	return bread;
}

void res_write(struct response* res, const char* buf, int len) {
	memcpy(&res->buf[res->len], buf, len);
	res->len += len;
}

void res_status(struct response* res, http_status status) {
	res->status = status;
}
