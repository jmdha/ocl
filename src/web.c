#include <assert.h>
#include <jhttp/jhttp.h>
#include <zstd.h>

#include "web.h"
#include "db.h"
#include "warlogs/warlogs.h"

struct route {
	char* method;
	char* path;
	int (*fn)(struct jhttp_response* res, const struct jhttp_request* req);
};

const char html_index[] = {
	#embed "../web/index.html"
};

const char css_style[] = {
	#embed "../web/style.css"
};

struct jhttp jhttp;

int get_index(struct jhttp_response* res, const struct jhttp_request* req) {
	snprintf(res->body, JHTTP_RESPONSE_MAX, "%s", html_index);
	res->status = 200;
	return 0;
}

int get_style(struct jhttp_response* res, const struct jhttp_request* req) {
	snprintf(res->body, JHTTP_RESPONSE_MAX, "%s", css_style);
	res->status = 200;
	return 0;
}

int post_api_upload(struct jhttp_response* res, const struct jhttp_request* req) {
	memset(res->body, 0, JHTTP_RESPONSE_MAX);
	int user_id;
	const char* file_name;
	for (size_t i = 0; i < req->header_count; i++)
		if (strcmp(req->headers[i].key, "File-Name") == 0)
			file_name = req->headers[i].val;
		else if (strcmp(req->headers[i].key, "User-ID") == 0)
			// TODO: improve safety
			user_id = atoi(req->headers[i].val);
	
	int log_id = db_logs_get_id(user_id, file_name);

	const char* ptr1;
	const char* ptr2;
	ptr1 = req->body;
	printf("body: \"%s\"\n", req->body);
	while ((ptr2 = strchr(ptr1, '\n')) != NULL) {
		int64_t ts;
		wl_event e;
		wl_error err;

		err = wl_parse(&ts, &e, ptr1, ptr2 - ptr1);
		if (err != wl_ok) {
			printf("failed to parse \"%.*s\"\n", (int) (ptr2 - ptr1), ptr1);
			break;
		}
		printf("%d\n", e.kind);
		ptr1 = ptr2;
	}

	res->status = 200;
	return 0;
}

const struct route routes[] = {
	{ "GET",  "/",           get_index },
	{ "GET",  "/index.html", get_index },
	{ "GET",  "/style.css",  get_style },
	{ "PUT",  "/api/upload", post_api_upload },
	{ "POST", "/api/upload", post_api_upload }
};

int handler(struct jhttp_response* res, const struct jhttp_request* req) {
	for (size_t i = 0; i < sizeof(routes) / sizeof(struct route); i++)
		if (strcmp(req->method, routes[i].method) == 0 && 
		    strcmp(req->path,   routes[i].path)   == 0) {
			db_requests_add(req->method, req->path);
			return routes[i].fn(res, req);
		}
	const char html_404[] = "<!DOCTYPE html><html><body><p>page not found</p></body></html>";
	res->status = 404;
	snprintf(res->body, JHTTP_RESPONSE_MAX, "%s", html_404);
	return 0;
}


void web_init() {
	if (jhttp_init(&jhttp, 8000, handler) != 0) {
		fprintf(stderr, "failed to initialize jhttp on port 8000\n");
		exit(1);
	}
}

void web_fini() {
	jhttp_fini(&jhttp);
}

void web_step() {
	jhttp_poll(&jhttp);
}
