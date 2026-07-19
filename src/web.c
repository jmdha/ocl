#include <mongoose/mongoose.h>
#include <assert.h>

#include "web.h"
#include "db.h"

struct mg_mgr mgr;

enum conn_kind {
	CONN_UNKNOWN,
	CONN_UPLOAD
};

struct upload_state {
	size_t expected;
	size_t received;
	void*  fp;
};
struct conn_state {
	enum conn_kind kind;
	union {
		struct upload_state us;
	};
};

static_assert(sizeof(struct conn_state) <= MG_DATA_SIZE, "conn_state too big");

void handle_metrics_requests(struct mg_connection* c, struct db* db) {
	int requests = db_metrics_requests(db);
	if (requests < 0) {
		mg_http_reply(c, 500, "", "%s", "internal server error\n");
		return;
	}
	mg_http_reply(c, 200, "", "%d", requests);
}

void handle_upload(struct mg_connection* c, int ev, void* ev_data) {
	struct conn_state* cs = (struct conn_state*) c->data;
	struct upload_state* us = &cs->us;
	struct mg_fs* fs = &mg_fs_posix;

	if (ev == MG_EV_HTTP_HDRS) {
		struct mg_http_message* hm = (struct mg_http_message*) ev_data;
		if (mg_match(hm->uri, mg_str("/api/upload"), NULL)) {
			cs->kind = CONN_UPLOAD;
			char buf[32];
			struct mg_str name = mg_str(mg_random_str(buf, sizeof(buf)));
			char path[96];
			mg_snprintf(path, sizeof(path), "%s/%.*s", "./uploads", name.len, name.buf);
			us->expected = hm->body.len;
			mg_iobuf_del(&c->recv, 0, hm->head.len);
			c->pfn = NULL;
			if (mg_path_is_sane(mg_str(path)))
				us->fp = fs->op(path, MG_FS_WRITE);
		}
	}

	if (cs->kind == CONN_UPLOAD && us->expected > 0 && c->recv.len > 0) {
		us->received += c->recv.len;
		if (us->fp) fs->wr(us->fp, c->recv.buf, c->recv.len);
		c->recv.len = 0;
		if (us->received >= us->expected) {
			MG_INFO(("uploaded %lu bytes", us->received));
			mg_http_reply(c, 200, NULL, "%lu ok\n", us->received);
			if (us->fp) fs->cl(us->fp);
			memset(cs, 0, sizeof(*cs));
			c->is_draining = 1;
		}
	}
}

void fn(struct mg_connection* c, int ev, void* ev_data) {
	struct db* db = (struct db*) c->fn_data;

	handle_upload(c, ev, ev_data);
	if (ev == MG_EV_HTTP_MSG && c->pfn != NULL) {
		struct mg_http_message* hm = (struct mg_http_message*) ev_data;
		db_metrics_requests_add(db, hm->method.len, hm->method.buf, hm->uri.len, hm->uri.buf);

		if (mg_match(hm->uri, mg_str("/api/metrics/requests"), NULL)) {
			handle_metrics_requests(c, db);
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
	mg_mgr_poll(&mgr, 1);
}
