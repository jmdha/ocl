#include "db.h"
#include "web.h"

int main(void) {
	struct db* db;
	db_init(&db);
	web_init(db);

	while(1)
		web_step();

	web_fini();
	db_fini(db);
}
