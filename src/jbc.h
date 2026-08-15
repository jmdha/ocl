#ifndef JBC_H
#define JBC_H

typedef struct jbc jbc;

int jbc_init(jbc** db, const char* path, size_t size, size_t max);
int jbc_fini(jbc* db);

size_t jbc_len(const jbc* db);
size_t jbc_size(const jbc* db);

int jbc_get(const jbc* db, size_t idx, void* data);
int jbc_ref(const jbc* db, size_t idx, const void** ref);
int jbc_add(jbc* db, size_t* idx, const void* data);

#endif
