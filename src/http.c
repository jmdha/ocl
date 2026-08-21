#include <openssl/sha.h>

#include "http.h"
#include "mongoose.h"
#include "utils.h"
#include "db.h"

void post_users(struct mg_connection* c, struct mg_http_message* hm) {
	db_user user;
	size_t  user_id;

	if (db_user_create(&user, &user_id) != 0)
		return mg_http_reply(c, 500, "", "");

	return mg_http_reply(c, 201, "", "%.*s", sizeof(user.key), user.key);
}

void post_login(struct mg_connection *c, struct mg_http_message *hm) {
	const db_user *user;
	size_t user_id;
	char token[256];

	mg_http_creds(hm, token, sizeof(token), token, sizeof(token));
	if (db_user_get_by_key(&user, &user_id, token) != 0)
		return mg_http_reply(c, 400, "", "Invalid Token\n");

	return mg_http_reply(
		c,
		303,
		"Location: /index.html\r\n"
		"Set-Cookie: access_token=%s; Path=/; HttpOnly; Secure; SameSite=Lax\r\n"
		"Content-Type: application/json\r\n",
		token
	);
}

void post_upload(struct mg_connection* c, struct mg_http_message* hm) {

	return mg_http_reply(c, 200, "", "");
}

void ev_handler(struct mg_connection* c, int ev, void* ev_data) {
	if (ev != MG_EV_HTTP_MSG) return;

	struct mg_http_message *hm = (struct mg_http_message *) ev_data;
	struct mg_http_serve_opts opts = {
		.root_dir = "web"
	};

	if (mg_match(hm->uri, mg_str("/api/public/users"), NULL)) {
		if (mg_match(hm->method, mg_str("POST"), NULL))
			post_users(c, hm);
		else
			mg_http_reply(c, 405, "", "Method Not Allowed\n");
	} else if (mg_match(hm->uri, mg_str("/api/public/login"), NULL)) {
		if (mg_match(hm->method, mg_str("POST"), NULL))
			post_login(c, hm);
		else
			mg_http_reply(c, 405, "", "Method Not Allowed\n");
	} else if (mg_match(hm->uri, mg_str("/api/upload"), NULL)) {
		if (mg_match(hm->method, mg_str("POST"), NULL))
			post_users(c, hm);
		else
			mg_http_reply(c, 405, "", "Method Not Allowed\n");
	} else {
		mg_http_serve_dir(c, hm, &opts);
	}
}
