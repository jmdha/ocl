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

UTEST(warlogs, parse_mapchange) {
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
