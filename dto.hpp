//
// Created by bishijie on 23-7-10.
//

#ifndef WEBBACK_DTO_HPP
#define WEBBACK_DTO_HPP
#include "oatpp/web/server/api/ApiController.hpp"
#include "oatpp/core/macro/codegen.hpp"
#include "oatpp/core/data/mapping/type/Object.hpp"
#include "oatpp/core/macro/codegen.hpp"
#include "oatpp/web/client/HttpRequestExecutor.hpp"
#include "oatpp-test/web/ClientServerTestRunner.hpp"
#include "oatpp/web/server/HttpConnectionHandler.hpp"
#include "oatpp/network/virtual_/client/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/server/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/Interface.hpp"
#include "oatpp/parser/json/mapping/ObjectMapper.hpp"
#include "oatpp/core/macro/component.hpp"
#include OATPP_CODEGEN_BEGIN(DTO)
using namespace oatpp::parser::json::mapping;
class user_dto : public oatpp::DTO {
    DTO_INIT(user_dto, DTO)
    DTO_FIELD(Int64 , id);
    DTO_FIELD(String, group);
    DTO_FIELD(String, name);
    DTO_FIELD(String, username);
    DTO_FIELD(String, password);
    DTO_FIELD(String, email);
};
class message_dto : public oatpp::DTO {
    DTO_INIT(message_dto, DTO);
    DTO_FIELD(String, message);
};
class message_error_dto : public oatpp::DTO {
    DTO_INIT(message_error_dto, DTO);
    DTO_FIELD(Int32, message);
};
#include OATPP_CODEGEN_END(DTO)
#endif //WEBBACK_DTO_HPP
