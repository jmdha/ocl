#include <pthread.h>
#include <unistd.h>

#include "worker.h"
#include "db.h"

static void worker_job(int job) {
}

static void* worker(void* arg) {
	int job;
	while (1) {
		if ((job = db_lock_job()) < 0) {
			sleep(1);
			continue;
		}
		worker_job(job);
	}
}

void workers(int count) {
	pthread_t threads[count];

	for (int i = 0; i < count; i++)
		pthread_create(&threads[i], NULL, worker, NULL);
}
