#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/socket.h>
#include <unistd.h>

#include "http.h"

#define BUF_MAX     (128 * 1024)
#define MAX_HEADERS 32
#define MAX_ROUTES  16

typedef struct {
	const char* method;
	const char* path;
	int (*f)(const char**, size_t, char*, size_t);
} route;

static route  routes[MAX_ROUTES];
static size_t route_count;

void http_register(const char* method, const char* path, int (*f)(const char**, size_t, char*, size_t)) {
	routes[route_count++] = (route) { method, path, f };
}

static const char* reason(int status) {
	switch (status) {
		case 200: return "OK";
		case 201: return "Created";
		case 303: return "See Other";
		case 400: return "Bad Request";
		case 401: return "Unauthorized";
		case 404: return "Not Found";
		case 409: return "Conflict";
		case 411: return "Length Required";
		default:  return "Internal Server Error";
	}
}

typedef struct {
	const char* method;
	const char* target;
	const char* hdrs[2 * MAX_HEADERS];
	size_t      count;
	char*       body;
	size_t      size;
} request;

// reads and parses the head, then reads the body; returns 0 on success, an http
// status on bad input, and negative when the peer is gone and no reply is due
static int request_read(int fd, char* buf, request* r) {
	size_t body_len = 0;
	size_t len      = 0;
	char*  blank;

	for (;;) {
		buf[len] = '\0';
		blank    = strstr(buf, "\r\n\r\n");
		if (blank)
			break;
		// a full buffer reads zero bytes, which drops the request like any other error
		ssize_t rd = read(fd, buf + len, BUF_MAX - len - 1);
		if (rd <= 0)
			return -1;
		len += (size_t) rd;
	}

	// status line: method, target, version
	char* nl = strstr(buf, "\r\n");
	*nl = '\0';
	r->method = strtok(buf, " ");
	r->target = strtok(NULL, " ");
	if (!r->method || !r->target || !strtok(NULL, " "))
		return 400;

	// header lines up to the blank one, cut into key/value pairs
	r->count = 0;
	for (char* line = nl + 2; line < blank;) {
		char* next = strstr(line, "\r\n");
		*next = '\0';

		char* colon = strchr(line, ':');
		if (!colon || r->count == MAX_HEADERS)
			return 400;
		*colon = '\0';

		r->hdrs[2 * r->count]     = line;
		r->hdrs[2 * r->count + 1] = colon + 1 + strspn(colon + 1, " \t");
		r->count++;
		line = next + 2;
	}

	for (size_t i = 0; i < r->count; i++) {
		if (strcasecmp(r->hdrs[2 * i], "Content-Length") == 0)
			body_len = strtoull(r->hdrs[2 * i + 1], NULL, 10);
		// no chunked decoding here, and a silently empty body would be acked as stored
		if (strcasecmp(r->hdrs[2 * i], "Transfer-Encoding") == 0)
			return 411;
	}

	size_t off = (size_t) (blank + 4 - buf);
	if (off + body_len >= BUF_MAX)
		return -1;

	while (len < off + body_len) {
		ssize_t rd = read(fd, buf + len, off + body_len - len);
		if (rd <= 0)
			return -1;
		len += (size_t) rd;
	}
	buf[off + body_len] = '\0';

	r->body = buf + off;
	r->size = BUF_MAX - off;
	return 0;
}

static int dispatch(request* r) {
	for (size_t i = 0; i < route_count; i++)
		if (strcmp(routes[i].method, r->method) == 0 && strcmp(routes[i].path, r->target) == 0)
			return routes[i].f(r->hdrs, r->count, r->body, r->size);
	return 404;
}

static void respond(int fd, int status, char* resp) {
	// error statuses answer with their reason phrase, whatever is left in buf is not a response
	if (status >= 400)
		resp = (char*) reason(status);

	// a handler may prepend its own header lines, separated from the body by a blank line
	const char* hd  = "";
	char*       sep = strstr(resp, "\r\n\r\n");
	if (sep) {
		sep[2] = '\0';
		hd     = resp;
		resp   = sep + 4;
	}

	dprintf(fd, "HTTP/1.1 %d %s\r\nConnection: close\r\n%sContent-Length: %zu\r\n\r\n%s",
	        status, reason(status), hd, strlen(resp), resp);
}

// one blocking request per connection: caddy fronts this, so no keep-alive and no pipelining
static void serve(int fd) {
	static char buf[BUF_MAX];
	request     req = { 0 };

	int status = request_read(fd, buf, &req);
	if (status < 0)
		return;
	if (status == 0)
		status = dispatch(&req);

	respond(fd, status, req.body);
}

void http_listen(int port) {
	// caddy fronts the api, so loopback only
	struct sockaddr_in sa = {
		.sin_family      = AF_INET,
		.sin_addr.s_addr = htonl(INADDR_LOOPBACK),
		.sin_port        = htons(port),
	};

	int listener = socket(AF_INET, SOCK_STREAM, 0);
	int opt      = 1;
	setsockopt(listener, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

	if (listener < 0 || bind(listener, (struct sockaddr*) &sa, sizeof(sa)) != 0 || listen(listener, SOMAXCONN) != 0) {
		perror("http_listen");
		exit(1);
	}

	for (;;) {
		int fd = accept(listener, NULL, NULL);
		if (fd < 0)
			continue;

		// a peer that stalls mid request would block the whole server, cut it loose instead
		struct timeval tv = { .tv_sec = 5 };
		setsockopt(fd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));

		serve(fd);
		close(fd);
	}
}
