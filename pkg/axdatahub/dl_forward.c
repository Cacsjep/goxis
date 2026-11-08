// Runtime loading of the Device Data Hub client (libdevice-data-hub-client-c).
//
// The library exists only on AXIS OS 12.10 and later. Linking it would stop the
// application from starting on any older camera ("error while loading shared
// libraries"), so the package does not link it. Each dh_ function the package
// calls is defined here, hidden, and forwards to the real one that
// axdatahub_load resolves with dlopen/dlsym. A forwarder called before a
// successful load does nothing and returns a zero value; the Go side calls
// Load first.
//
// When the package starts calling a new dh_ function, declare it in
// datahub.h, add a forwarder here and resolve it in axdatahub_resolve, or the
// link fails with an undefined reference.

#include <dlfcn.h>
#include <pthread.h>
#include <stddef.h>
#include <stdio.h>

// Hidden: the executable must not export these names, or the library's own
// calls to them would bind to the forwarders.
#pragma GCC visibility push(hidden)
#include "datahub.h"

#define AXDATAHUB_FWD(ret, name, params, args, fallback)                     \
	static __typeof__(&name) p_##name;                                       \
	ret name params { return p_##name ? p_##name args : fallback; }

#define AXDATAHUB_FWD_VOID(name, params, args)                               \
	static __typeof__(&name) p_##name;                                       \
	void name params { if (p_##name) p_##name args; }

AXDATAHUB_FWD(const char *, dh_error_get_message, (const DHError *err), (err), NULL)
AXDATAHUB_FWD(DHErrorCode, dh_error_get_code, (const DHError *err), (err), 0)
AXDATAHUB_FWD(const char *, dh_error_code_to_string, (DHErrorCode code), (code), NULL)
AXDATAHUB_FWD_VOID(dh_error_destroy, (DHError *err), (err))
AXDATAHUB_FWD(int64_t, dh_timestamp_get_as_ms, (const DHTimestamp *timestamp), (timestamp), 0)
AXDATAHUB_FWD(const char *, dh_topic_data_get_json_data, (const DHTopicData *topic_data), (topic_data), NULL)
AXDATAHUB_FWD(const char *, dh_topic_sample_get_topic_name, (const DHTopicSample *topic_sample), (topic_sample), NULL)
AXDATAHUB_FWD(uint64_t, dh_topic_sample_get_sample_id, (const DHTopicSample *topic_sample), (topic_sample), 0)
AXDATAHUB_FWD(bool, dh_topic_sample_is_historical, (const DHTopicSample *topic_sample), (topic_sample), false)
AXDATAHUB_FWD(const DHTimestamp *, dh_topic_sample_get_timestamp, (const DHTopicSample *topic_sample), (topic_sample), NULL)
AXDATAHUB_FWD(const DHTopicData *, dh_topic_sample_get_data, (const DHTopicSample *topic_sample), (topic_sample), NULL)
AXDATAHUB_FWD(DHInstanceKeys *, dh_instance_keys_create, (void), (), NULL)
AXDATAHUB_FWD_VOID(dh_instance_keys_destroy, (DHInstanceKeys *keys), (keys))
AXDATAHUB_FWD(bool, dh_instance_keys_add_integer, (DHInstanceKeys *keys, const char *name, int64_t value, DHError **error), (keys, name, value, error), false)
AXDATAHUB_FWD(uint32_t, dh_topic_list_get_count, (const DHTopicList *list), (list), 0)
AXDATAHUB_FWD(const char *, dh_topic_list_get_name, (const DHTopicList *list, uint32_t index), (list, index), NULL)
AXDATAHUB_FWD_VOID(dh_topic_list_destroy, (DHTopicList *topic_list), (topic_list))
AXDATAHUB_FWD(DHClient *, dh_client_create, (const char *name, DHError **error), (name, error), NULL)
AXDATAHUB_FWD(bool, dh_client_set_connection_update_callback, (DHClient *client, DHOnConnectionUpdateCallback callback, void *user_data, DHError **error), (client, callback, user_data, error), false)
AXDATAHUB_FWD(bool, dh_client_connect, (DHClient *client, DHError **error), (client, error), false)
AXDATAHUB_FWD_VOID(dh_client_destroy, (DHClient *client), (client))
AXDATAHUB_FWD(DHTopicList *, dh_client_get_topic_list, (DHClient *client, DHError **error), (client, error), NULL)
AXDATAHUB_FWD(DHSubscriber *, dh_client_create_subscriber, (DHClient *client, const char *subscriber_name, DHError **error), (client, subscriber_name, error), NULL)
AXDATAHUB_FWD(bool, dh_client_set_logging, (DHClient *client, DHLogLevel level, DHLogTarget target, DHError **error), (client, level, target, error), false)
AXDATAHUB_FWD(DHFilter *, dh_filter_create, (void), (), NULL)
AXDATAHUB_FWD_VOID(dh_filter_destroy, (DHFilter *filter), (filter))
AXDATAHUB_FWD(bool, dh_filter_add_topic_name, (DHFilter *filter, const char *topic_name, DHError **error), (filter, topic_name, error), false)
AXDATAHUB_FWD(bool, dh_filter_add_instance, (DHFilter *filter, const DHInstanceKeys *instance_keys, DHError **error), (filter, instance_keys, error), false)
AXDATAHUB_FWD(DHSubscribeOptions *, dh_subscribe_options_create, (void), (), NULL)
AXDATAHUB_FWD_VOID(dh_subscribe_options_destroy, (DHSubscribeOptions *options), (options))
AXDATAHUB_FWD(bool, dh_subscribe_options_add_filter, (DHSubscribeOptions *options, const DHFilter *filter, DHError **error), (options, filter, error), false)
AXDATAHUB_FWD_VOID(dh_subscribe_options_set_enable_data_updates, (DHSubscribeOptions *options, bool enable), (options, enable))
AXDATAHUB_FWD_VOID(dh_subscriber_destroy, (DHSubscriber *subscriber), (subscriber))
AXDATAHUB_FWD(bool, dh_subscriber_set_data_callback, (DHSubscriber *subscriber, DHOnDataCallback on_data, void *user_data, DHError **error), (subscriber, on_data, user_data, error), false)
AXDATAHUB_FWD(bool, dh_subscriber_subscribe, (DHSubscriber *subscriber, const DHSubscribeOptions *options, DHError **error), (subscriber, options, error), false)

#pragma GCC visibility pop

#include "axdatahub.h"

// Trampolines with the exact callback types the library expects. The user data
// is a cgo.Handle; the Go exports take it back as uintptr_t, and take a
// non-const sample.
static void axdatahub_connection_trampoline(DHConnectionState update, void *user_data) {
	axdatahubOnConnectionUpdate(update, (uintptr_t)user_data);
}

static void axdatahub_data_trampoline(const DHTopicSample *topic_sample, void *user_data) {
	axdatahubOnData((DHTopicSample *)topic_sample, (uintptr_t)user_data);
}

bool axdatahub_set_connection_handle(DHClient *client, uintptr_t handle, DHError **error) {
	return dh_client_set_connection_update_callback(client, axdatahub_connection_trampoline, (void *)handle, error);
}

bool axdatahub_set_data_handle(DHSubscriber *subscriber, uintptr_t handle, DHError **error) {
	return dh_subscriber_set_data_callback(subscriber, axdatahub_data_trampoline, (void *)handle, error);
}

static int axdatahub_state = 0; // 0 not tried, 1 loaded, -1 failed
static char axdatahub_error[256];
static pthread_mutex_t axdatahub_load_lock = PTHREAD_MUTEX_INITIALIZER;

#define AXDATAHUB_RESOLVE(lib, name)                                         \
	do {                                                                     \
		*(void **)(&p_##name) = dlsym(lib, #name);                           \
		if (p_##name == NULL) {                                              \
			snprintf(axdatahub_error, sizeof(axdatahub_error),               \
			         "missing symbol %s", #name);                            \
			return -1;                                                       \
		}                                                                    \
	} while (0)

#define AXDATAHUB_SYMBOLS(X)                                                 \
	X(dh_error_get_message)                                                  \
	X(dh_error_get_code)                                                     \
	X(dh_error_code_to_string)                                               \
	X(dh_error_destroy)                                                      \
	X(dh_timestamp_get_as_ms)                                                \
	X(dh_topic_data_get_json_data)                                           \
	X(dh_topic_sample_get_topic_name)                                        \
	X(dh_topic_sample_get_sample_id)                                         \
	X(dh_topic_sample_is_historical)                                         \
	X(dh_topic_sample_get_timestamp)                                         \
	X(dh_topic_sample_get_data)                                              \
	X(dh_instance_keys_create)                                               \
	X(dh_instance_keys_destroy)                                              \
	X(dh_instance_keys_add_integer)                                          \
	X(dh_topic_list_get_count)                                               \
	X(dh_topic_list_get_name)                                                \
	X(dh_topic_list_destroy)                                                 \
	X(dh_client_create)                                                      \
	X(dh_client_set_connection_update_callback)                              \
	X(dh_client_connect)                                                     \
	X(dh_client_destroy)                                                     \
	X(dh_client_get_topic_list)                                              \
	X(dh_client_create_subscriber)                                           \
	X(dh_client_set_logging)                                                 \
	X(dh_filter_create)                                                      \
	X(dh_filter_destroy)                                                     \
	X(dh_filter_add_topic_name)                                              \
	X(dh_filter_add_instance)                                                \
	X(dh_subscribe_options_create)                                           \
	X(dh_subscribe_options_destroy)                                          \
	X(dh_subscribe_options_add_filter)                                       \
	X(dh_subscribe_options_set_enable_data_updates)                          \
	X(dh_subscriber_destroy)                                                 \
	X(dh_subscriber_set_data_callback)                                       \
	X(dh_subscriber_subscribe)

static int axdatahub_resolve(void *lib) {
#define AXDATAHUB_RESOLVE_ONE(name) AXDATAHUB_RESOLVE(lib, name);
	AXDATAHUB_SYMBOLS(AXDATAHUB_RESOLVE_ONE)
#undef AXDATAHUB_RESOLVE_ONE
	return 0;
}

// A library that lacks one symbol is not used at all: no forwarder may reach
// a half-resolved table.
static void axdatahub_clear(void) {
#define AXDATAHUB_CLEAR_ONE(name) p_##name = NULL;
	AXDATAHUB_SYMBOLS(AXDATAHUB_CLEAR_ONE)
#undef AXDATAHUB_CLEAR_ONE
}

// axdatahub_load opens the library once. It returns 1 when every forwarder is
// resolved, else 0 with the reason in axdatahub_load_error. A failed load is
// final: the library does not appear at runtime.
int axdatahub_load(void) {
	pthread_mutex_lock(&axdatahub_load_lock);
	if (axdatahub_state == 0) {
		axdatahub_state = -1;
		void *lib = dlopen("libdevice-data-hub-client-c.so.0", RTLD_NOW | RTLD_LOCAL);
		if (lib == NULL) {
			const char *e = dlerror();
			snprintf(axdatahub_error, sizeof(axdatahub_error), "%s", e ? e : "dlopen failed");
		} else if (axdatahub_resolve(lib) == 0) {
			axdatahub_state = 1;
		} else {
			axdatahub_clear();
			dlclose(lib);
		}
	}
	int ok = axdatahub_state == 1;
	pthread_mutex_unlock(&axdatahub_load_lock);
	return ok;
}

const char *axdatahub_load_error(void) { return axdatahub_error; }
