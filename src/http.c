#include "http.h"
#include "mongoose.h"
#include "utils.h"
#include "db.h"

void post_users(struct mg_connection* c, struct mg_http_message* hm) {
	db_user user;
	size_t user_id;

	if (gen_ran(user.key, sizeof(user.key)) != 0)
		return mg_http_reply(c, 500, "", "");

	if (db_user_add(&user, &user_id) != 0)
		return mg_http_reply(c, 500, "", "");

	return mg_http_reply(c, 201, "", "%.*s", sizeof(user.key), user.key);
}

void post_login(struct mg_connection *c, struct mg_http_message *hm) {
	size_t user_id;
	char key[32];
	char buf[512];

	mg_http_creds(hm, key, sizeof(key), key, sizeof(key));
	if (db_user_get_id(&user_id, key) != 0)
		return mg_http_reply(c, 401, "", "");

	snprintf(
		buf,
		sizeof(buf),
		"Location: /index.html\r\n"
		"Set-Cookie: access_token=%s; Path=/; HttpOnly; Secure; SameSite=Lax\r\n"
		"Content-Type: application/json\r\n",
		key
	);

	return mg_http_reply(
		c,
		303,
		buf,
		""
	);
}

void get_upload(struct mg_connection* c, struct mg_http_message* hm) {
	db_log log;
	size_t user_id;
	size_t log_id;
	struct mg_str* fn_header;
	char filename[32];
	char key[32];
	char buf[64];
	struct stat st;

	mg_http_creds(hm, key, sizeof(key), key, sizeof(key));
	if (db_user_get_id(&user_id, key) != 0)
		return mg_http_reply(c, 401, "", "");

	if (!(fn_header = mg_http_get_header(hm, "Filename")))
		return mg_http_reply(c, 400, "", "Missing Filename\n");

	if (fn_header->len >= sizeof(filename))
		return mg_http_reply(c, 400, "", "Filename Too Long\n");

	memcpy(filename, fn_header->buf, fn_header->len);
	filename[fn_header->len] = '\0';

	if (db_log_get_id(&log_id, user_id, filename) != 0) {
		log.user_id = user_id;
		strcpy(log.filename, filename);
		if (db_log_add(&log, &log_id) != 0)
			return mg_http_reply(c, 500, "", "");
	}

	snprintf(buf, sizeof(buf), "logs/%zu", log_id);
	if (stat(buf, &st) != 0)
		return mg_http_reply(c, 200, "", "%zu", 0);
	else
		return mg_http_reply(c, 200, "", "%zu", st.st_size);
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
		if (mg_match(hm->method, mg_str("GET"), NULL))
			get_upload(c, hm);
		else if (mg_match(hm->method, mg_str("POST"), NULL))
			post_upload(c, hm);
		else
			mg_http_reply(c, 405, "", "Method Not Allowed\n");
	} else {
		mg_http_serve_dir(c, hm, &opts);
	}
}
