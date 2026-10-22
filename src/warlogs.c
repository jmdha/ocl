#include <string.h>
#include <stdio.h>

#include "warlogs.h"

int64_t wl_timestamp(int d, int m, int y, int H, int M, int S, int f) {
	return (int64_t)(y - 1970) * 31622400000LL +
	       (int64_t)(m - 1)    * 2678400000LL +
	       (int64_t)(d - 1)    * 86400000LL +
	       (int64_t)(H)        * 3600000LL +
	       (int64_t)(M)        * 60000LL +
	       (int64_t)(S)        * 1000LL +
	       (int64_t)(f);
}

bool wl_parse_version(wl_event* e, const char* str) {
	return sscanf(
		str,
		"COMBAT_LOG_VERSION,%d,ADVANCED_LOG_ENABLED,%d,BUILD_VERSION,%d.%d.%d,PROJECT_ID,%d",
		&e->version.log_version,
		&e->version.advanced,
		&e->version.major,
		&e->version.minor,
		&e->version.patch,
		&e->version.project
	) == 6;
}

bool wl_parse_zonechange(wl_event* e, const char* str) {
	return sscanf(
		str,
		"ZONE_CHANGE,%d,\"%[^\"]\",%d",
		&e->zonechange.instance_id,
		&e->zonechange.zone_name,
		&e->zonechange.difficulty_id
	) == 3;
}

bool wl_parse_mapchange(wl_event* e, const char* str) {
	return sscanf(
		str,
		"MAP_CHANGE,%d,\"%[^\"]\",%f,%f,%f,%f",
		&e->mapchange.map_id,
		&e->mapchange.map_name,
		&e->mapchange.x0,
		&e->mapchange.x1,
		&e->mapchange.y0,
		&e->mapchange.y1
	) == 6;
}

wl_error wl_parse(int64_t* ts, wl_event* e, const char* str) {
	int d, m, y, H, M, S, f;
	int pos;
	int r;
	r = sscanf(str, "%d/%d/%d %d:%d:%d.%d  %n", &m, &d, &y, &H, &M, &S, &f, &pos);
	if (r != 7)
		return wl_invalid_format;
	if (d < 1 || d > 31 || m < 1 || m > 12 || y < 2000 || y > 2040)
		return wl_invalid_date;
	if (H < 0 || H > 23 || M < 0 || M > 59 || S < 0 || S > 60)
		return wl_invalid_time;
	str += pos;
	*ts = wl_timestamp(d, m, y, H, M, S, f);

	struct {
		wl_kind kind;
		bool (*f)(wl_event*, const char*);
	} cases[] = {
		{ wl_kind_version,    wl_parse_version },
		{ wl_kind_zonechange, wl_parse_zonechange },
		{ wl_kind_mapchange,  wl_parse_mapchange }
	};

	for (size_t i = 0; i < sizeof(cases) / sizeof(cases[0]); i++)
		if (cases[i].f(e, str)) {
			e->kind = cases[i].kind;
			return wl_ok;
		}

	e->kind = wl_kind_unknown;
	return wl_invalid;
}
