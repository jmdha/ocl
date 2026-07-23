#ifndef NET_H
#define NET_H

#include "db.h"

void web_init(struct db* db);
void web_fini();
void web_step();

#endif
