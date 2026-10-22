#include <utest.h/utest.h>

#include "warlogs.h"

UTEST(warlogs, timestamp) {
	ASSERT_EQ(0,           wl_timestamp(1, 1, 1970, 0, 0, 0, 0));
	ASSERT_EQ(86400000,    wl_timestamp(2, 1, 1970, 0, 0, 0, 0));
	ASSERT_EQ(2678400000,  wl_timestamp(1, 2, 1970, 0, 0, 0, 0));
	ASSERT_EQ(31622400000, wl_timestamp(1, 1, 1971, 0, 0, 0, 0));
	ASSERT_EQ(3600000,     wl_timestamp(1, 1, 1970, 1, 0, 0, 0));
	ASSERT_EQ(60000,       wl_timestamp(1, 1, 1970, 0, 1, 0, 0));
	ASSERT_EQ(1000,        wl_timestamp(1, 1, 1970, 0, 0, 1, 0));
	ASSERT_EQ(1,           wl_timestamp(1, 1, 1970, 0, 0, 0, 1));
}

UTEST(warlogs, parse_version) {
	int64_t ts;
	wl_event e;
	wl_error err;

	err = wl_parse(&ts, &e, "6/5/2026 16:26:02.0002  COMBAT_LOG_VERSION,22,ADVANCED_LOG_ENABLED,1,BUILD_VERSION,12.0.5,PROJECT_ID,1");

	ASSERT_EQ(wl_ok,           err);
	ASSERT_EQ(wl_kind_version, e.kind);
	ASSERT_EQ(22,              e.version.log_version);
	ASSERT_EQ(1,               e.version.advanced);
	ASSERT_EQ(12,              e.version.major);
	ASSERT_EQ(0,               e.version.minor);
	ASSERT_EQ(5,               e.version.patch);
	ASSERT_EQ(1,               e.version.project);
}

UTEST(warlogs, parse_zonechange) {
	int64_t ts;
	wl_event e;
	wl_error err;

	err = wl_parse(&ts, &e, "9/16/2026 19:36:34.0502  ZONE_CHANGE,2993,\"Vaults of Atal'Utek\",23");

	ASSERT_EQ(wl_ok,                    err);
	ASSERT_EQ(wl_kind_zonechange,       e.kind);
	ASSERT_EQ(2993,                     e.zonechange.instance_id);
	ASSERT_STREQ("Vaults of Atal'Utek", e.zonechange.zone_name);
	ASSERT_EQ(23,                       e.zonechange.difficulty_id);
}

UTEST(warlogs, parse_mapchange) {
	int64_t ts;
	wl_event e;
	wl_error err;

	err = wl_parse(&ts, &e, "9/16/2026 19:36:34.0502  MAP_CHANGE,2588,\"Altar of Fangs\",2008.329956,1414.579956,2552.080078,1662.500000");

	ASSERT_EQ(wl_ok,               err);
	ASSERT_EQ(wl_kind_mapchange,   e.kind);
	ASSERT_EQ(2588,                e.mapchange.map_id);
	ASSERT_STREQ("Altar of Fangs", e.mapchange.map_name);
	ASSERT_NEAR(2008.329956,         e.mapchange.x0, 0.01f);
	ASSERT_NEAR(1414.579956,         e.mapchange.x1, 0.01f);
	ASSERT_NEAR(2552.080078,         e.mapchange.y0, 0.01f);
	ASSERT_NEAR(1662.500000,         e.mapchange.y1, 0.01f);
}
