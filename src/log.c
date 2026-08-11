#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <time.h>

#include "log.h"

#define RED     "\x1b[31m"
#define GREEN   "\x1b[32m"
#define YELLOW  "\x1b[33m"
#define BLUE    "\x1b[34m"
#define MAGENTA "\x1b[35m"
#define CYAN    "\x1b[36m"
#define GRAY    "\x1b[90m"
#define RESET   "\x1b[0m"

void log_init(void) {
}

void log_out(int level, const char* msg, ...) {
	char formatted[32000];
	
	va_list arg_ptr;
	va_start(arg_ptr, msg);
	vsnprintf(formatted, sizeof(formatted), msg, arg_ptr);
	va_end(arg_ptr);
	
	// Get current UTC time
	time_t now = time(NULL);
	struct tm *utc = gmtime(&now);
	
	char out_time[64];
	
	snprintf(out_time, sizeof(out_time),
	         "[%04d-%02d-%02dT%02d:%02d:%02dZ]",
	         utc->tm_year + 1900,
	         utc->tm_mon + 1,
	         utc->tm_mday,
	         utc->tm_hour,
	         utc->tm_min,
	         utc->tm_sec);
	
	if (level == LOG_LEVEL_TRACE)
		fprintf(stdout, GRAY "%s %s" RESET "\n", out_time, formatted);
	else
		fprintf(stdout, "%s %s\n", out_time, formatted);
}
