#ifndef AXDATAHUB_PRIVATE_H_INCLUDED
#define AXDATAHUB_PRIVATE_H_INCLUDED

// Package-private declarations: the Go exports, the loader and the callback
// registration helpers. Kept out of datahub.h so dl_forward.c can declare the
// dh_ forwarders hidden without hiding these.

#include "datahub.h"

// Go callbacks (datahub.go). The user data is a runtime/cgo.Handle.
extern void axdatahubOnConnectionUpdate(DHConnectionState update, uintptr_t handle);
extern void axdatahubOnData(DHTopicSample *topic_sample, uintptr_t handle);

// dl_forward.c: runtime loading of libdevice-data-hub-client-c.
extern int axdatahub_load(void);
extern const char *axdatahub_load_error(void);

// Register the Go callbacks with a cgo.Handle as user data (dl_forward.c). The
// handle travels as uintptr_t so Go never converts it to a pointer.
bool axdatahub_set_connection_handle(DHClient *client, uintptr_t handle, DHError **error);
bool axdatahub_set_data_handle(DHSubscriber *subscriber, uintptr_t handle, DHError **error);

#endif
