#include <signal.h>
#include <stdio.h>
#include <cjson/cJSON.h>
#include <jhttp/jhttp.h>

#include "db.h"
#include "log.h"
#include "utils.h"

static int sig_exit = 0;
static void signal_handler(int sig) {
	printf("signal %d received\n", sig);
	sig_exit = 1;
}

int post_users(struct jhttp_response* res, const struct jhttp_request* req) {
	db_user user;
	char key[64];

	if (db_user_new(&user, key) != 0) {
		snprintf(res->body, sizeof(res->body), "failed to create user");
		res->status = 400;
		return 0;
	}

	snprintf(res->body, sizeof(res->body), "%s", key);
	res->status = 200;
	return 0;
}

int post_login(struct jhttp_response* res, const struct jhttp_request* req) {
	db_user user;

	if (db_user_get_by_key(&user, req->body) != 0) {
		snprintf(res->body, sizeof(res->body), "unknown key");
		res->status = 400;
		return 0;
	}

	snprintf(res->body, sizeof(res->body), "%lu", user.id);
	res->status = 200;
	return 0;
}

int get_upload(struct jhttp_response* res, const struct jhttp_request* req) {
	return 0;
	//const char* filename = NULL;
	//printf("get_upload\n");
	//for (size_t i = 0; i < req->header_count; i++)
	//	if (strcmp(req->headers[i].key, "Filename") == 0)
	//		filename = req->headers[i].val;
	//if (!filename) {
	//	printf("no filename\n");
	//	snprintf(res->body, JHTTP_RESPONSE_MAX, "");
	//	res->status = 400;
	//	return 0;
	//}
	//int line = db_file_get(filename);
	//if (line < 0) {
	//	printf("negative line\n");
	//	snprintf(res->body, JHTTP_RESPONSE_MAX, "");
	//	res->status = 400;
	//	return 0;
	//}
	//printf("setting res 200\n");
	//snprintf(res->body, JHTTP_RESPONSE_MAX, "%d", line);
	//res->status = 200;
	//return 0;
}

int post_upload(struct jhttp_response* res, const struct jhttp_request* req) {
	return 0;
	//int line;
	//char filename[64];
	//char content[4192];
	//printf("printing len\n");
	//printf("len: %s\n", strlen(req->body));
	//if (sscanf(req->body, "%s %d %4191[^\n]", filename, &line, content) != 3) {
	//	printf("failed to parse body: \"%s\"\n", req->body);
	//	snprintf(res->body, JHTTP_RESPONSE_MAX, "");
	//	res->status = 400;
	//	return 0;
	//}
	//if (db_file_set(filename, line, content) != 0) {
	//	snprintf(res->body, JHTTP_RESPONSE_MAX, "");
	//	res->status = 400;
	//	return 0;
	//}
	//res->status = 200;
	//snprintf(res->body, JHTTP_RESPONSE_MAX, "");
	//return 0;
}

typedef struct {
	const char* method;
	const char* path;
	int (*func)(struct jhttp_response* res, const struct jhttp_request* req);
} route;

const route routes[] = {
	{ "POST", "/users", post_users },
	{ "POST", "/login", post_login },
	{ "GET",  "/upload", get_upload },
	{ "POST", "/upload", post_upload },
};

int handler(struct jhttp_response* res, const struct jhttp_request* req) {
	INFO("%s: %s", req->method, req->path);
	for (size_t i = 0; i < sizeof(routes) / sizeof(routes[0]); i++)
		if (strcmp(req->method, routes[i].method) == 0 &&
		    strcmp(req->path, routes[i].path) == 0)
			return routes[i].func(res, req);

	const char html_404[] = "<!DOCTYPE html><html><body><p>page not found</p></body></html>";
	res->status = 404;
	snprintf(res->body, sizeof(res->body), "%s", html_404);
	return 0;
}

int main(void) {
	struct jhttp jhttp;
	jhttp_init(&jhttp, 8000, handler);
	db_init();

	signal(SIGINT,  signal_handler);
	signal(SIGTERM, signal_handler);
	signal(SIGKILL, signal_handler);
	signal(SIGQUIT, signal_handler);
	signal(SIGABRT, signal_handler);

	while(sig_exit == 0)
		jhttp_poll(&jhttp);

	db_fini();
	jhttp_fini(&jhttp);
}
