#ifndef HTTP_H
#define HTTP_H

#include <stddef.h>

// registers route for http listener; the handler gets key/value header pairs and the request
// body in buf, leaves its response in buf, and returns the status to send back
void http_register(const char* method, const char* path, int (*f)(const char**, size_t, char*, size_t));
// listen on port and serve requests forever; the kernel closes the sockets when the process dies
void http_listen(int port);

#endif
