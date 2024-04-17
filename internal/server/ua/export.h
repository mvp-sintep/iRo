#ifndef _INC_OPCUA_H_
#define _INC_OPCUA_H_

#include "open62541.h"

typedef UA_StatusCode ( *callback_t )(
  UA_Server *server,
  const UA_NodeId *sessionId,
  void *sessionContext,
  const UA_NodeId *nodeId,
  void *nodeContext,
  UA_Boolean includeSourceTimeStamp,
  const UA_NumericRange *range,
  UA_DataValue *value );

UA_Server *New( uint32_t port, uint32_t parentNodes_length, int namespaceId );
int CreateObjectNode( UA_Server *server, char *nodeID, char *nodeName );
int CreateI32DataSource( UA_Server *server, char *nodeID, char *nodeName, char *parentNodeID, int32_t defaultValue );
int CreateI16DataSource( UA_Server *server, char *nodeID, char *nodeName, char *parentNodeID, int16_t defaultValue );
int CreateFloatDataSource( UA_Server *server, char *nodeID, char *nodeName, char *parentNodeID, float defaultValue );
void Run( UA_Server *server );

#endif // _INC_OPCUA_H_
