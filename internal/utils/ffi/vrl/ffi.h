#pragma once

#include <stdlib.h>


typedef struct {
	char* key;
	char* value;
} hmap_entry;

typedef struct {
	size_t      count;
	hmap_entry* entries;
} hmap;

typedef void* script_runner;


extern void hmap_free(hmap* self);

extern script_runner script_runner_new(const char* source, char** err);
extern void          script_runner_free(script_runner self);
extern hmap*         script_runner_eval(script_runner runner, hmap* input, char** err);
