#ifndef LOG_H
#define LOG_H

#define LOG_LEVEL_INFO 1
#define LOG_LEVEL_TRACE 2

void log_init(void);
void log_out(int level, const char *msg, ...);

#if LOG_INFO
#define INFO(msg, ...) log_out(LOG_LEVEL_INFO, msg, ##__VA_ARGS__);
#else
#define INFO(msg, ...)
#endif

#if LOG_TRACE
#define TRACE(msg, ...) log_out(LOG_LEVEL_TRACE, msg, ##__VA_ARGS__);
#else
#define TRACE(msg, ...)
#endif

#endif
