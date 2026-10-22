#ifndef WARLOGS_H
#define WARLOGS_H

#include <stddef.h>
#include <stdint.h>

typedef enum {
	wl_ok,
	wl_invalid,
	wl_invalid_format,
	wl_invalid_date,
	wl_invalid_time
} wl_error;

typedef enum {
	wl_kind_version,
	wl_kind_zonechange,
	wl_kind_mapchange,
	wl_kind_unknown,
} wl_kind;

typedef struct {
	int log_version;
	int advanced;
	int major;
	int minor;
	int patch;
	int project;
} wl_version;

typedef struct {
	int  instance_id;
	int  difficulty_id;
	char zone_name[64];
} wl_zonechange;

typedef struct {
	int   map_id;
	char  map_name[64];
	float x0, x1, y0, y1;
} wl_mapchange;

typedef struct {
	wl_kind kind;
	union {
		wl_version    version;
		wl_zonechange zonechange;
		wl_mapchange  mapchange;
	};
} wl_event;

int64_t wl_timestamp(int d, int m, int y, int H, int M, int S, int f);
wl_error wl_parse(int64_t* ts, wl_event* e, const char* str);

#endif
