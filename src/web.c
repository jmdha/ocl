#include <mongoose/mongoose.h>

#include "web.h"
#include "db.h"

struct mg_mgr mgr;

void handle_metrics_requests(struct mg_connection* c, struct db* db) {
	int requests = db_metrics_requests(db);
	if (requests < 0) {
		mg_http_reply(c, 500, "", "%s", "internal server error\n");
		return;
	}
	mg_http_reply(c, 200, "", "%d", requests);
}

void handle_upload(struct mg_connection* c) {

}

void fn(struct mg_connection* c, int ev, void* ev_data) {
	struct db* db = (struct db*) c->fn_data;

	if (ev == MG_EV_HTTP_MSG) {
		struct mg_http_message* hm = (struct mg_http_message*) ev_data;
		db_metrics_requests_add(db, hm->method.len, hm->method.buf, hm->uri.len, hm->uri.buf);

		if (mg_match(hm->uri, mg_str("/api/metrics/requests"), NULL)) {
			handle_metrics_requests(c, db);
		} else if (mg_match(hm->uri, mg_str("/api/upload"), NULL)) {
			handle_upload(c);
		} else {
			struct mg_http_serve_opts opts;
			memset(&opts, 0, sizeof(opts));
			opts.root_dir = "./web";
			mg_http_serve_dir(c, ev_data, &opts);
		}

		MG_INFO((
			"%lu %.*s %.*s -> %.*s",
			c->id,
			(int) hm->method.len,
			hm->method.buf,
			(int) hm->uri.len,
			hm->uri.buf,
			(int) 3,
			&c->send.buf[9]
		));
	}
}

void web_init(struct db* db) {
	mg_mgr_init(&mgr);
	mg_http_listen(&mgr, "localhost:8000", fn, db);
}

void web_fini() {
	mg_mgr_free(&mgr);
}

void web_step() {
	mg_mgr_poll(&mgr, 0);
}
