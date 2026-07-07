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
#include "worker.h"
#include "log.h"

static char* log_dir = "upload";

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
		log_error("r_api_upload(db_add_upload): error");
		res_status(res, STATUS_INTERNAL_SERVER_ERROR);
		return;
	}

	if (ranfile(&fp, path, sizeof(path), log_dir) != 0) {
		log_error("r_api_upload(ranfile): error");
		db_set_upload_error(id, "ranfile", strlen("ranfile"));
		res_status(res, STATUS_INTERNAL_SERVER_ERROR);
		return;
	}

	if (db_set_upload_path(id, path, strlen(path)) != 0) {
		log_error("r_api_upload(path): error");
		db_set_upload_error(id, "upload path", strlen("upload path"));
		res_status(res, STATUS_INTERNAL_SERVER_ERROR);
		return;
	}

	size = 0;
	while ((r = req_read(req, buf, sizeof(buf))) != 0) {
		if (r < 0) {
			log_error("r_api_upload(pipe): error");
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

	workers(4);

	server_add("/",               &r_index);
	server_add("/index.html",     &r_index);
	server_add("/api/characters", &r_api_characters);
	server_add("/api/upload",     &r_api_upload);
	server_listen(8082);
}
