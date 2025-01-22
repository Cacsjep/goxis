#ifndef MDB_H_INCLUDED
#define MDB_H_INCLUDED

#include <mdb/message.h>
#include <mdb/connection.h>
#include <mdb/error.h>
#include <mdb/subscriber.h>

// Only declarations here:
extern void onConnectionErrorCallback(mdb_error_t *error, void *user_data);
extern void onMessageCallback(mdb_message_t *message, void *user_data);
extern void onSubscriberCreateDoneCallback(mdb_error_t *error, void *user_data);

#endif