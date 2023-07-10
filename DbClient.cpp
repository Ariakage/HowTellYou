//
// Created by bishijie on 23-7-10.
//
#include "oatpp/orm/SchemaMigration.hpp"
#include "oatpp/orm/DbClient.hpp"
#include "dto.hpp"
#include "oatpp/core/macro/codegen.hpp"
#include "oatpp/web/server/HttpConnectionHandler.hpp"
#include "oatpp/network/virtual_/client/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/server/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/Interface.hpp"
#include "oatpp/parser/json/mapping/ObjectMapper.hpp"
#include "oatpp/core/macro/component.hpp"
#include "oatpp/parser/json/mapping/ObjectMapper.hpp"
#include "oatpp/web/server/HttpConnectionHandler.hpp"
#include "oatpp/network/tcp/server/ConnectionProvider.hpp"
#include "oatpp/core/macro/component.hpp"
#include "oatpp/orm/SchemaMigration.hpp"
#include "oatpp/orm/DbClient.hpp"
#include "oatpp/core/data/stream/BufferStream.hpp"
#include "oatpp/core/macro/codegen.hpp"
#include OATPP_CODEGEN_BEGIN(DbClient)
class user_client : public oatpp::orm::DbClient {
public:
    user_client(const std::shared_ptr<oatpp::orm::Executor>& executor)
            : oatpp::orm::DbClient(executor)
    {}
    QUERY(login,
          "SELECT * FROM users WHERE username=:username;",
          PARAM(oatpp::String, username))
    QUERY(addf,
          "SELECT * FROM users WHERE username=:username;",
          PARAM(oatpp::String, username))
    QUERY(reg,
          "SELECT * FROM users WHERE username=:username;",
          PARAM(oatpp::String, username))
};
#include OATPP_CODEGEN_END(DbClient)