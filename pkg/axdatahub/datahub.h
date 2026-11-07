#ifndef AXDATAHUB_DATAHUB_H_INCLUDED
#define AXDATAHUB_DATAHUB_H_INCLUDED

// The subset of the Device Data Hub C client API (libdevice-data-hub-client-c,
// AXIS OS 12.10 and later) that this package calls.
//
// The package loads the library at runtime (dl_forward.c), so it must also
// build with SDKs that ship no Device Data Hub headers, such as ACAP Native
// SDK 1.15 for AXIS OS 11 and 12. The declarations below are therefore written
// out here. When the SDK does ship the real headers, they are included first
// and every declaration below is checked against them by the compiler: a
// changed signature is a compile error, not a silent ABI mismatch.

#include <stdbool.h>
#include <stdint.h>

#if defined(__has_include)
#if __has_include(<datahub/client.h>)
#include <datahub/client.h>
#define AXDATAHUB_HAVE_SDK_HEADERS 1
#endif
#endif

typedef uint32_t DHErrorCode;
typedef uint32_t DHConnectionState;
typedef uint32_t DHLogLevel;
typedef uint32_t DHLogTarget;

typedef struct DHError_t DHError;
typedef struct DHClient_t DHClient;
typedef struct DHSubscriber_t DHSubscriber;
typedef struct DHFilter_t DHFilter;
typedef struct DHSubscribeOptions_t DHSubscribeOptions;
typedef struct DHInstanceKeys_t DHInstanceKeys;
typedef struct DHTopicSample_t DHTopicSample;
typedef struct DHTopicData_t DHTopicData;
typedef struct DHTimestamp_t DHTimestamp;
typedef struct DHTopicList_t DHTopicList;

typedef void (*DHOnConnectionUpdateCallback)(DHConnectionState update, void *user_data);
typedef void (*DHOnDataCallback)(const DHTopicSample *topic_sample, void *user_data);

// Values from datahub/client.h and datahub/common.h (SDK 12.11).
#define AXDATAHUB_CONN_DISCONNECTED ((DHConnectionState)0)
#define AXDATAHUB_CONN_CONNECTED ((DHConnectionState)1)
#define AXDATAHUB_LOG_WARNING ((DHLogLevel)3)
#define AXDATAHUB_LOG_TARGET_SYSLOG ((DHLogTarget)1)

#if defined(AXDATAHUB_HAVE_SDK_HEADERS)
_Static_assert(AXDATAHUB_CONN_DISCONNECTED == DH_CONN_DISCONNECTED, "DH_CONN_DISCONNECTED changed");
_Static_assert(AXDATAHUB_CONN_CONNECTED == DH_CONN_CONNECTED, "DH_CONN_CONNECTED changed");
_Static_assert(AXDATAHUB_LOG_WARNING == DH_LOG_WARNING, "DH_LOG_WARNING changed");
_Static_assert(AXDATAHUB_LOG_TARGET_SYSLOG == DH_LOG_TARGET_SYSLOG, "DH_LOG_TARGET_SYSLOG changed");
#endif

// common.h
const char *dh_error_get_message(const DHError *err);
DHErrorCode dh_error_get_code(const DHError *err);
const char *dh_error_code_to_string(DHErrorCode code);
void dh_error_destroy(DHError *err);
int64_t dh_timestamp_get_as_ms(const DHTimestamp *timestamp);
const char *dh_topic_data_get_json_data(const DHTopicData *topic_data);
const char *dh_topic_sample_get_topic_name(const DHTopicSample *topic_sample);
uint64_t dh_topic_sample_get_sample_id(const DHTopicSample *topic_sample);
bool dh_topic_sample_is_historical(const DHTopicSample *topic_sample);
const DHTimestamp *dh_topic_sample_get_timestamp(const DHTopicSample *topic_sample);
const DHTopicData *dh_topic_sample_get_data(const DHTopicSample *topic_sample);
DHInstanceKeys *dh_instance_keys_create(void);
void dh_instance_keys_destroy(DHInstanceKeys *keys);
bool dh_instance_keys_add_integer(DHInstanceKeys *keys, const char *name, int64_t value, DHError **error);
uint32_t dh_topic_list_get_count(const DHTopicList *list);
const char *dh_topic_list_get_name(const DHTopicList *list, uint32_t index);
void dh_topic_list_destroy(DHTopicList *topic_list);

// client.h
DHClient *dh_client_create(const char *name, DHError **error);
bool dh_client_set_connection_update_callback(DHClient *client, DHOnConnectionUpdateCallback callback,
                                              void *user_data, DHError **error);
bool dh_client_connect(DHClient *client, DHError **error);
void dh_client_destroy(DHClient *client);
DHTopicList *dh_client_get_topic_list(DHClient *client, DHError **error);
DHSubscriber *dh_client_create_subscriber(DHClient *client, const char *subscriber_name, DHError **error);
bool dh_client_set_logging(DHClient *client, DHLogLevel level, DHLogTarget target, DHError **error);

// subscriber.h
DHFilter *dh_filter_create(void);
void dh_filter_destroy(DHFilter *filter);
bool dh_filter_add_topic_name(DHFilter *filter, const char *topic_name, DHError **error);
bool dh_filter_add_instance(DHFilter *filter, const DHInstanceKeys *instance_keys, DHError **error);
DHSubscribeOptions *dh_subscribe_options_create(void);
void dh_subscribe_options_destroy(DHSubscribeOptions *options);
bool dh_subscribe_options_add_filter(DHSubscribeOptions *options, const DHFilter *filter, DHError **error);
void dh_subscribe_options_set_enable_data_updates(DHSubscribeOptions *options, bool enable);
void dh_subscriber_destroy(DHSubscriber *subscriber);
bool dh_subscriber_set_data_callback(DHSubscriber *subscriber, DHOnDataCallback on_data, void *user_data,
                                     DHError **error);
bool dh_subscriber_subscribe(DHSubscriber *subscriber, const DHSubscribeOptions *options, DHError **error);

#endif
