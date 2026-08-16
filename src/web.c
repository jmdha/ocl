#include <openssl/sha.h>

#include "web.h"
#include "mongoose.h"
#include "utils.h"
#include "db.h"

struct mg_mgr mgr;

void post_users(struct mg_connection* c, struct mg_http_message* hm) {
	db_user user;
	size_t user_id;

	if (gen_ran(user.key, sizeof(user.key)) != 0) {
		mg_http_reply(c, 500, "", "Internal server error\n");
		return;
	}

	if (db_user_add(&user, &user_id) != 0) {
		mg_http_reply(c, 500, "", "Internal server error\n");
		return;
	}

	mg_http_reply(c, 201, "", "%.*s", sizeof(user.key), user.key);
}

void post_login(struct mg_connection *c, struct mg_http_message *hm) {
	const db_user* user;
	size_t user_id;
	char token[256];
	char headers[320];

	int len = mg_http_get_var(&hm->body, "token", token, sizeof(token));

	if (len <= 0) {
		mg_http_reply(c, 400, "", "Missing token\n");
		return;
	}

	if (db_user_get_by_key(&user, &user_id, token) != 0) {
		mg_http_reply(c, 400, "", "Invalid token\n");
		return;
	}

	mg_snprintf(headers, sizeof(headers),
		    "Location: /index.html\r\n"
	            "Set-Cookie: access_token=%s; Path=/; HttpOnly; Secure; SameSite=Lax\r\n"
	            "Content-Type: application/json\r\n",
	            token);

	mg_http_reply(c, 303, headers, "{}\n");
}

void post_upload(struct mg_connection* c, struct mg_http_message* hm, size_t user_id) {
	mg_http_reply(c, 501, "", "");

}

int auth(struct mg_http_message* hm, size_t* user_id) {
	struct mg_str* cookie = mg_http_get_header(hm, "Cookie");

	if (cookie == NULL) {
		printf("no cookie\n");
		return 1;
	}

	struct mg_str token = mg_http_get_header_var(*cookie, mg_str("access_token"));
	if (!token.buf || token.len == 0) {
		printf("invalid token str\n");
		return 1;
	}

	const db_user* user;
	if (db_user_get_by_key(&user, user_id, token.buf) != 0) {
		return 1;
	}

	return 0;
}

void ev_handler(struct mg_connection* c, int ev, void* ev_data) {
	if (ev != MG_EV_HTTP_MSG) return;

	struct mg_http_message *hm = (struct mg_http_message *) ev_data;

	if (mg_match(hm->uri, mg_str("/api/public/#"), NULL)) {
		if (mg_match(hm->method, mg_str("POST"), NULL) && mg_match(hm->uri, mg_str("/api/public/users"), NULL))
			post_users(c, hm);
		else if (mg_match(hm->method, mg_str("POST"), NULL) && mg_match(hm->uri, mg_str("/api/public/login"), NULL))
			post_login(c, hm);
		else
			mg_http_reply(c, 404, "", "Not found\n");
	} else if (mg_match(hm->uri, mg_str("/api/#"), NULL)) {
		size_t user_id;
		if (auth(hm, &user_id) != 0) {
			mg_http_reply(c, 401, "", "Unauthorized\n");
			return;
		}

		if (mg_match(hm->method, mg_str("POST"), NULL) && mg_match(hm->uri, mg_str("/api/upload"), NULL))
			post_users(c, hm);
		else
			mg_http_reply(c, 404, "", "Not found\n");
	} else {
		struct mg_http_serve_opts opts = {
			.root_dir = "web"
		};
		mg_http_serve_dir(c, hm, &opts);
	}
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
