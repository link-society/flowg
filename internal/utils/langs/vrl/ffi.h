#pragma once

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>


typedef void* script_runner;

typedef struct msgpack_buffer msgpack_buffer;
struct msgpack_buffer {
	const uint8_t* data;
	size_t         len;
};

typedef struct compilation_result compilation_result;
struct compilation_result {
	bool is_ok;

	struct {
		script_runner runner;
	} ok;

	struct {
		char* reason;
	} err;
};

typedef struct log_transform_result log_transform_result;
struct log_transform_result {
	bool is_ok;

	struct {
		msgpack_buffer data;
	} ok;

	struct {
		char* reason;
	} err;
};

extern compilation_result compile_script    (const char* source);
extern void               drop_script_runner(script_runner self);

extern log_transform_result transform_log(script_runner runner, msgpack_buffer log_event);
