#include "export.h"

static int NamespaceID = 1;

UA_Server *New( uint32_t port, uint32_t parentNodes_length, int namespaceId )
{
  NamespaceID = namespaceId;
  UA_Server *server = UA_Server_new();
  UA_ServerConfig_setMinimal( UA_Server_getConfig( server ), port, NULL );
  UA_Server_getConfig( server )->verifyRequestTimestamp = UA_RULEHANDLING_ACCEPT;
  return server;
}

int CreateObjectNode( UA_Server *server, char *nodeID, char *nodeName )
{
  UA_NodeId ua_parentNodeId = UA_NODEID_NUMERIC( 0, UA_NS0ID_OBJECTSFOLDER );
  UA_NodeId ua_parentReferenceTypeId = UA_NODEID_NUMERIC( 0, UA_NS0ID_ORGANIZES );
  UA_QualifiedName ua_browseName = UA_QUALIFIEDNAME( 1, nodeName );
  UA_NodeId ua_nodeTypeDefinition = UA_NODEID_NUMERIC( 0, UA_NS0ID_BASEOBJECTTYPE );
  UA_NodeId ua_nodeID = UA_NODEID_STRING( 1, nodeID );
  UA_ObjectAttributes attr = UA_ObjectAttributes_default;
  attr.displayName = UA_LOCALIZEDTEXT( "en-US", nodeName );

  return UA_Server_addObjectNode(
    server,
    ua_nodeID,
    ua_parentNodeId,
    ua_parentReferenceTypeId,
    ua_browseName,
    ua_nodeTypeDefinition,
    attr,
    NULL,
    NULL );
}

extern int16_t go_read_callback_int32( const int16_t address, int32_t *value );
extern int16_t go_read_callback_int16( const int16_t address, int16_t *value );
extern int16_t go_read_callback_float( const int16_t address, float *value );

extern int16_t go_write_callback_int32( const int16_t address, int32_t *value );
extern int16_t go_write_callback_int16( const int16_t address, int16_t *value );
extern int16_t go_write_callback_float( const int16_t address, float *value );

static UA_StatusCode readInt32CallBack(
    UA_Server *server,
    const UA_NodeId *sessionId,
    void *sessionContext,
    const UA_NodeId *nodeId,
    void *nodeContext,
    UA_Boolean sourceTimeStamp,
    const UA_NumericRange *range,
    UA_DataValue *dataValue )
{
  int32_t tmp = 0;
  if( nodeId->identifier.string.data != NULL ) {
    if( go_read_callback_int32( atoi( nodeId->identifier.string.data + 3 ) * 4, &tmp ) != 0 ) {
      UA_Variant_setScalarCopy( &dataValue->value, &tmp, &UA_TYPES[ UA_TYPES_INT32 ] );
      dataValue->hasValue = true;
      return UA_STATUSCODE_GOOD;
    }
  }
  return UA_STATUSCODE_BAD;
}

static UA_StatusCode readInt16CallBack(
    UA_Server *server,
    const UA_NodeId *sessionId,
    void *sessionContext,
    const UA_NodeId *nodeId,
    void *nodeContext,
    UA_Boolean sourceTimeStamp,
    const UA_NumericRange *range,
    UA_DataValue *dataValue )
{
  int16_t tmp = 0;
  if( nodeId->identifier.string.data != NULL ) {
    if( go_read_callback_int16( atoi( nodeId->identifier.string.data + 3 ) * 2, &tmp ) != 0 ) {
      UA_Variant_setScalarCopy( &dataValue->value, &tmp, &UA_TYPES[ UA_TYPES_INT16 ] );
      dataValue->hasValue = true;
      return UA_STATUSCODE_GOOD;
    }
  }
  return UA_STATUSCODE_BAD;
}

static UA_StatusCode readFloatCallBack(
    UA_Server *server,
    const UA_NodeId *sessionId,
    void *sessionContext,
    const UA_NodeId *nodeId,
    void *nodeContext,
    UA_Boolean sourceTimeStamp,
    const UA_NumericRange *range,
    UA_DataValue *dataValue )
{
  float tmp = 0;
  if( nodeId->identifier.string.data != NULL ) {
    if( go_read_callback_float( atoi( nodeId->identifier.string.data + 1 ) * 4, &tmp ) != 0 ) {
      UA_Variant_setScalarCopy( &dataValue->value, &tmp, &UA_TYPES[ UA_TYPES_FLOAT ] );
      dataValue->hasValue = true;
      return UA_STATUSCODE_GOOD;
    }
  }
  return UA_STATUSCODE_BAD;
}

static UA_StatusCode writeInt32CallBack(
  UA_Server *server,
  const UA_NodeId *sessionId,
  void *sessionContext,
  const UA_NodeId *nodeId,
  void *nodeContext,
  const UA_NumericRange *range,
  const UA_DataValue *value )
{
  int32_t tmp = 0;
  if( nodeId->identifier.string.data != NULL && value->hasValue ) {
    tmp = *((int32_t*)(value->value.data));
    if( go_write_callback_int32( atoi( nodeId->identifier.string.data + 3 ) * 4, &tmp ) != 0 ) {
      return UA_STATUSCODE_GOOD;
    }
  }
  return UA_STATUSCODE_BAD;
}

static UA_StatusCode writeInt16CallBack(
  UA_Server *server,
  const UA_NodeId *sessionId,
  void *sessionContext,
  const UA_NodeId *nodeId,
  void *nodeContext,
  const UA_NumericRange *range,
  const UA_DataValue *value )
{
  int16_t tmp = 0;
  if( nodeId->identifier.string.data != NULL && value->hasValue ) {
    tmp = *((int16_t*)(value->value.data));
    if( go_write_callback_int16( atoi( nodeId->identifier.string.data + 3 ) * 2, &tmp ) != 0 ) {
      return UA_STATUSCODE_GOOD;
    }
  }
  return UA_STATUSCODE_BAD;
}

static UA_StatusCode writeFloatCallBack(
  UA_Server *server,
  const UA_NodeId *sessionId,
  void *sessionContext,
  const UA_NodeId *nodeId,
  void *nodeContext,
  const UA_NumericRange *range,
  const UA_DataValue *value )
{
  float tmp = 0;
  if( nodeId->identifier.string.data != NULL && value->hasValue ) {
    tmp = *((float*)(value->value.data));
    if( go_write_callback_float( atoi( nodeId->identifier.string.data + 1 ) * 4, &tmp ) != 0 ) {
      return UA_STATUSCODE_GOOD;
    }
  }
  return UA_STATUSCODE_BAD;
}

int createDataSourceTag(
    UA_Server *server,
    char *nodeID,
    char *nodeName,
    char *parentNodeID,
    UA_Variant defaultValue,
    uint32_t ua_type,
    read_callback_t read_callback,
    write_callback_t write_callback)
{
  UA_NodeId ua_parentNodeId = UA_NODEID_STRING( 1, parentNodeID );
  UA_NodeId ua_parentReferenceNodeId = UA_NODEID_NUMERIC( 0, UA_NS0ID_ORGANIZES );

  UA_NodeId ua_nodeID = UA_NODEID_STRING( NamespaceID, nodeID );
  UA_NodeId ua_nodeTypeDefinition = UA_NODEID_NUMERIC( 0, UA_NS0ID_BASEDATAVARIABLETYPE );
  UA_QualifiedName ua_nodeName = UA_QUALIFIEDNAME( NamespaceID, nodeName );

  UA_VariableAttributes ua_attributes = UA_VariableAttributes_default;
  ua_attributes.displayName = UA_LOCALIZEDTEXT( "en-US", nodeName );
  ua_attributes.dataType = UA_TYPES[ ua_type ].typeId;
  ua_attributes.accessLevel = UA_ACCESSLEVELMASK_READ | UA_ACCESSLEVELMASK_WRITE;
  ua_attributes.value = defaultValue;

  UA_DataSource ua_dataSource;
  ua_dataSource.read = read_callback;
  ua_dataSource.write = write_callback;

  return UA_Server_addDataSourceVariableNode(
    server,
    ua_nodeID,
    ua_parentNodeId,
    ua_parentReferenceNodeId,
    ua_nodeName,
    ua_nodeTypeDefinition,
    ua_attributes,
    ua_dataSource,
    NULL,
    NULL);
}

int CreateI32DataSource(
    UA_Server *server,
    char *nodeID,
    char *nodeName,
    char *parentNodeID,
    int32_t defaultValue )
{
  int ua_type = UA_TYPES_INT32;
  UA_Variant ua_defaultValue;
  UA_Variant_init( &ua_defaultValue );
  UA_Variant_setScalar( &ua_defaultValue, &defaultValue, &UA_TYPES[ ua_type ] );

  return createDataSourceTag(
    server,
    nodeID,
    nodeName,
    parentNodeID,
    ua_defaultValue,
    ua_type,
    readInt32CallBack,
    writeInt32CallBack );
}

int CreateI16DataSource(
    UA_Server *server,
    char *nodeID,
    char *nodeName,
    char *parentNodeID,
    int16_t defaultValue )
{
  int ua_type = UA_TYPES_INT16;
  UA_Variant ua_defaultValue;
  UA_Variant_init( &ua_defaultValue );
  UA_Variant_setScalar( &ua_defaultValue, &defaultValue, &UA_TYPES[ ua_type ] );

  return createDataSourceTag(
    server,
    nodeID,
    nodeName,
    parentNodeID,
    ua_defaultValue,
    ua_type,
    readInt16CallBack,
    writeInt16CallBack );
}

int CreateFloatDataSource(
    UA_Server *server,
    char *nodeID,
    char *nodeName,
    char *parentNodeID,
    float defaultValue )
{
  int ua_type = UA_TYPES_FLOAT;
  UA_Variant ua_defaultValue;
  UA_Variant_init( &ua_defaultValue );
  UA_Variant_setScalar( &ua_defaultValue, &defaultValue, &UA_TYPES[ ua_type ] );

  return createDataSourceTag(
    server,
    nodeID,
    nodeName,
    parentNodeID,
    ua_defaultValue,
    ua_type,
    readFloatCallBack,
    writeFloatCallBack );
}

void Run( UA_Server *server )
{
  bool isServerRunning = true;
  UA_Server_run( server, &isServerRunning );
  UA_Server_delete( server );
}
