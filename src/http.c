#include <openssl/sha.h>

#include "http.h"
#include "mongoose.h"
#include "utils.h"
#include "db.h"

void users(struct mg_connection* c, struct mg_http_message* hm) {
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

void login(struct mg_connection *c, struct mg_http_message *hm) {
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

void upload(struct mg_connection* c, struct mg_http_message* hm) {
	mg_http_reply(c, 501, "", "");

}

void ev_handler(struct mg_connection* c, int ev, void* ev_data) {
	if (ev != MG_EV_HTTP_MSG) return;

	struct mg_http_message *hm = (struct mg_http_message *) ev_data;
	struct mg_http_serve_opts opts = {
		.root_dir = "web"
	};

	if (mg_match(hm->uri, mg_str("/api/public/users"), NULL))
		users(c, hm);
	else if (mg_match(hm->uri, mg_str("/api/public/login"), NULL))
		login(c, hm);
	else if (mg_match(hm->uri, mg_str("/api/upload"), NULL))
		upload(c, hm);
	else
		mg_http_serve_dir(c, hm, &opts);
}
