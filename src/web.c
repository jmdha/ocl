#include <openssl/sha.h>

#include "web.h"
#include "mongoose.h"
#include "utils.h"
#include "db.h"

struct mg_mgr mgr;

uint32_t post_users(struct mg_connection* c, struct mg_http_message* hm) {
	db_user user;
	size_t user_id;

	if (gen_ran(user.key, sizeof(user.key)) != 0) {
		mg_http_reply(c, 500, "", "Internal server error\n");
		return 500;
	}
	if (db_user_add(&user, &user_id) != 0) {
		mg_http_reply(c, 500, "", "Internal server error\n");
		return 500;
	}

	mg_http_reply(c, 201, "", "%s\n", user.key);
	return 201;
}

uint32_t post_login(struct mg_connection *c, struct mg_http_message *hm) {
	const db_user* user;
	size_t user_id;
	char token[256];
	char headers[320];

	int len = mg_http_get_var(&hm->body, "token", token, sizeof(token));

	if (len <= 0) {
		mg_http_reply(c, 400, "", "Missing token\n");
		return 400;
	}
	if (db_user_get_by_key(&user, &user_id, token) != 0) {
		mg_http_reply(c, 400, "", "Invalid token\n");
		return 400;
	}

	mg_snprintf(headers, sizeof(headers),
		    "Location: /index.html\r\n"
	            "Set-Cookie: access_token=%s; Path=/; HttpOnly; Secure; SameSite=Lax\r\n"
	            "Content-Type: application/json\r\n",
	            token);

	mg_http_reply(c, 303, headers, "{}\n");
	return 303;
}

void ev_handler(struct mg_connection* c, int ev, void* ev_data) {
	if (ev != MG_EV_HTTP_MSG) return;

	struct mg_http_message *hm = (struct mg_http_message *) ev_data;
	db_request req = { 0 };
	size_t req_id;

	snprintf(req.path,   sizeof(req.path),   "%.*s", (int) hm->uri.len,    hm->uri.buf);
	snprintf(req.method, sizeof(req.method), "%.*s", (int) hm->method.len, hm->method.buf);
	memcpy(req.ip, c->rem.addr.ip, sizeof(req.ip));
	req.is_ip6 = c->rem.is_ip6;
	req.ts = current_time_ns();

	if (mg_match(hm->uri, mg_str("/api/#"), NULL)) {
		if (mg_match(hm->method, mg_str("POST"), NULL) && mg_match(hm->uri, mg_str("/api/users"), NULL))
			req.status = post_users(c, hm);
		else if (mg_match(hm->method, mg_str("POST"), NULL) && mg_match(hm->uri, mg_str("/api/login"), NULL))
			req.status = post_login(c, hm);
		else {
			req.status = 404;
			mg_http_reply(c, 404, "", "Not found\n");
		}
	} else {
		struct mg_http_serve_opts opts = {
			.root_dir = "web"
		};
		mg_http_serve_dir(c, hm, &opts);
	}

	req.dur = current_time_ns() - req.ts;
	db_request_add(&req, &req_id);
}

int web_init(int port) {
	char addr[128];
	snprintf(addr, sizeof(addr), "http://0.0.0.0:%d", port);
	
	mg_mgr_init(&mgr);
	
	struct mg_connection *c = mg_http_listen(&mgr, addr, ev_handler, NULL);
	if (c == NULL) {
		mg_mgr_free(&mgr);
		return -1;
	}
	
	return 0;
}

int web_fini() {
	mg_mgr_free(&mgr);
	return 0;
}

int web_poll() {
	mg_mgr_poll(&mgr, 1);
	return 0;
}
