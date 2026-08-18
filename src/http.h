#ifndef HTTP_H
#define HTTP_H

#include "mongoose.h"

void ev_handler(struct mg_connection* c, int ev, void* ev_data);

#endif
