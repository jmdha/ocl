#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <unistd.h>

#include "routes.h"

#include "db.h"
#include "utils.h"

static const char* header_get(const char* headers[], size_t count, const char* key) {
	// header names are case insensitive
	for (size_t i = 0; i < count; i++)
		if (strcasecmp(headers[2 * i], key) == 0)
			return headers[2 * i + 1];
	return NULL;
}

static int form_get(const char* body, const char* key, char* out, size_t len) {
	char   tmp[512];
	char   pat[64];
	size_t pat_len = (size_t) snprintf(pat, sizeof(pat), "%s=", key);

	if (strlen(body) >= sizeof(tmp))
		return 1;
	strcpy(tmp, body);

	for (char* p = strtok(tmp, "&"); p; p = strtok(NULL, "&")) {
		if (strncmp(p, pat, pat_len) != 0 || strlen(p + pat_len) >= len)
			continue;
		strcpy(out, p + pat_len);
		return 0;
	}

	return 1;
}

// the access_token cookie the login route hands out is the only credential; returns an http status on failure
static int log_open(const char* headers[], size_t count, size_t* log_id, char* path, size_t len) {
	db_log log;
	char   cookie[512];

	const char* val = header_get(headers, count, "cookie");
	if (!val)
		return 401;
	snprintf(cookie, sizeof(cookie), "%s", val);

	char* key = strstr(cookie, "access_token=");
	if (!key)
		return 401;

	key += strlen("access_token=");
	key[strcspn(key, ";")] = '\0';

	if (db_user_get_id(&log.user_id, key) != 0)
		return 401;

	const char* filename = header_get(headers, count, "filename");
	if (!filename || strchr(filename, '/') || strlen(filename) >= sizeof(log.filename))
		return 400;
	strcpy(log.filename, filename);

	if (db_log_get_id(log_id, log.user_id, log.filename) != 0 && db_log_add(&log, log_id) != 0)
		return 500;

	snprintf(path, len, "logs/%s", log.filename);
	printf("upload \"%s\"\n", log.filename);
	return 0;
}

int post_users(const char* headers[], size_t count, char* buf, size_t size) {
	(void) headers;
	(void) count;
	db_user user;
	size_t  user_id;

	if (gen_ran(user.key, sizeof(user.key)) != 0 || db_user_add(&user, &user_id) != 0)
		return 500;

	snprintf(buf, size, "%s", user.key);
	return 201;
}

int post_login(const char* headers[], size_t count, char* buf, size_t size) {
	(void) headers;
	(void) count;
	size_t user_id;
	char   key[32];

	if (form_get(buf, "token", key, sizeof(key)) != 0 || db_user_get_id(&user_id, key) != 0)
		return 401;

	snprintf(buf, size,
	         "Location: /index.html\r\n"
	         "Set-Cookie: access_token=%s; Path=/; HttpOnly; Secure; SameSite=Lax\r\n"
	         "\r\n",
	         key);
	return 303;
}

// tells the client how much of the log we already hold, so it can resume
int get_upload(const char* headers[], size_t count, char* buf, size_t size) {
	size_t log_id;
	char   path[128];

	int err = log_open(headers, count, &log_id, path, sizeof(path));
	if (err)
		return err;

	snprintf(buf, size, "%zu", file_size(path));
	return 200;
}

int post_upload(const char* headers[], size_t count, char* buf, size_t size) {
	(void) size;
	size_t log_id;
	char   path[128];
	FILE*  fp;

	int err = log_open(headers, count, &log_id, path, sizeof(path));
	if (err)
		return err;

	const char* obuf = header_get(headers, count, "offset");
	if (!obuf)
		return 400;

	// chunks must land in order, right where the last one ended
	size_t offset = strtoull(obuf, NULL, 10);
	if (offset != file_size(path))
		return 409;

	fp = fopen(path, "a");
	if (!fp)
		return 500;

	// a half written chunk would leave the resume offset mid line, so roll it back
	size_t body_len = strlen(buf);
	if (fwrite(buf, 1, body_len, fp) != body_len || fclose(fp) != 0) {
		truncate(path, (off_t) offset);
		return 500;
	}

	buf[0] = '\0';
	return 200;
}
