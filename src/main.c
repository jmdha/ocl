#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <cjson/cJSON.h>
#include <unistd.h>

#include "server.h"
#include "db.h"
#include "embed.h"
#include "utils.h"

static char* log_dir = "upload";

void r_style(struct request* req, struct response* res) {
	res_write(res, STYLE, sizeof(STYLE));
	res_status(res, STATUS_OK);
}

void r_index(struct request* req, struct response* res) {
	res_write(res, INDEX, sizeof(INDEX));
	res_status(res, STATUS_OK);
}

void r_api_characters(struct request* req, struct response* res) {

}

void r_api_upload(struct request* req, struct response* res) {
	size_t r, size;
	char   buf[64 * 1024];
	char   path[1024];
	FILE*  fp;
	int    id;

	if ((id = db_add_upload()) < 0) {
		printf("db_add_upload failed\n");
		res_status(res, STATUS_INTERNAL_SERVER_ERROR);
		return;
	}

	if (ranfile(&fp, path, sizeof(path), log_dir) != 0) {
		printf("ranfile failed\n");
		db_set_upload_error(id, "ranfile", strlen("ranfile"));
		res_status(res, STATUS_INTERNAL_SERVER_ERROR);
		return;
	}

	if (db_set_upload_path(id, path, strlen(path)) != 0) {
		printf("db_set_upload_path failed\n");
		db_set_upload_error(id, "upload path", strlen("upload path"));
		res_status(res, STATUS_INTERNAL_SERVER_ERROR);
		return;
	}

	size = 0;
	while ((r = req_read(req, buf, sizeof(buf))) != 0) {
		if (r < 0) {
			printf("pipe read failed\n");
			fclose(fp);
			unlink(path);
			db_set_upload_error(id, "buffer read", strlen("buffer read"));
			return;
		}
		fwrite(buf, r, 1, fp);
		size += r;
	}
	fclose(fp);

	db_set_upload_done(id, size);

	res_status(res, STATUS_OK);
}

int main(int argc, char** argv) {
	signal(SIGPIPE, SIG_IGN);
	srand(time(NULL));

	db_init("ocl.db");

	server_add("/style.css",      &r_style);
	server_add("/",               &r_index);
	server_add("/index.html",     &r_index);
	server_add("/api/characters", &r_api_characters);
	server_add("/api/upload",     &r_api_upload);
	server_listen(8082);
}
