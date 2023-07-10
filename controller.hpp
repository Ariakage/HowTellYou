//
// Created by bishijie on 23-7-10.
//

#ifndef WEBBACK_CONTROLLER_HPP
#define WEBBACK_CONTROLLER_HPP
#include "oatpp/web/server/api/ApiController.hpp"
#include "oatpp/core/macro/codegen.hpp"
#include "oatpp/web/client/HttpRequestExecutor.hpp"
#include "oatpp-test/web/ClientServerTestRunner.hpp"
#include "oatpp/web/server/HttpConnectionHandler.hpp"
#include "oatpp/network/virtual_/client/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/server/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/Interface.hpp"
#include "oatpp/parser/json/mapping/ObjectMapper.hpp"
#include "oatpp/core/macro/component.hpp"
#include "dto.hpp"
#include "oatpp/orm/SchemaMigration.hpp"
#include "oatpp/orm/DbClient.hpp"
#include "oatpp/core/macro/codegen.hpp"
#include "oatpp/core/data/stream/BufferStream.hpp"
#include "DbClient.cpp"
#include OATPP_CODEGEN_BEGIN(ApiController)
oatpp::Object<user_dto> temp,temp2,temp3,temp4;
oatpp::Object<message_dto> me;
oatpp::Object<message_error_dto> mee;
class user_controller : public oatpp::web::server::api::ApiController {
public:
    user_controller(OATPP_COMPONENT(std::shared_ptr<ObjectMapper>, objectMapper) /* Inject object mapper */)
            : oatpp::web::server::api::ApiController(objectMapper)
    {}
    ENDPOINT("GET", "/", root,BODY_STRING(String, start)) {
        OATPP_ASSERT_HTTP(start, Status::CODE_404, "error");
        mee->message = 114514;
        return createResponse(Status::CODE_200, "OK");
    }
    ENDPOINT_ASYNC("POST", "/reg",userim) {
    ENDPOINT_ASYNC_INIT(userim)
        Action act() override {
            return request->readBodyToDtoAsync<oatpp::Object<user_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&userim::returnResponse);
        }
        Action returnResponse(const oatpp::Object<user_dto>& body){
            OATPP_ASSERT_HTTP(body, Status::CODE_404, "error");
            /* check user information
             * save body information
             * code........
             * */
            return _return(controller->createResponse(Status::CODE_100,"OK"));
        }
    };
    ENDPOINT_ASYNC("POST", "/login",userim2) {
    ENDPOINT_ASYNC_INIT(userim2)
        Action act() override {
            return request->readBodyToDtoAsync<oatpp::Object<user_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&userim::returnResponse);
        }
        Action returnResponse(const oatpp::Object<user_dto>& body){
            OATPP_ASSERT_HTTP(body, Status::CODE_404, "error");
            return _return(controller->createDtoResponse(Status::CODE_100,me));
        }
    };
    ENDPOINT_ASYNC("POST", "/getms",gms) {
    ENDPOINT_ASYNC_INIT(gms)
        Action act() override {
            return request->readBodyToDtoAsync<oatpp::Object<user_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&gms::use);
        }
        Action use(const oatpp::Object<user_dto>& body){
            temp = body;
        }
        Action use2(const oatpp::Object<user_dto>& body){
            bool b1 = *temp->name==*temp2->name;
            bool b2 = *temp->password==*temp2->password;
            bool b3 = *temp->username==*temp2->username;
            bool b4 = *temp->email==*temp2->email;
            bool b5 = *temp->id==*temp2->id;
            if(b1&&b2&&b3&&b4&&b5){
                return _return(controller->createDtoResponse(Status::CODE_100,me));
            }
            return _return(controller->createDtoResponse(Status::CODE_400,mee));
        }
    };
    ENDPOINT_ASYNC("POST", "/sendim",sim) {
    ENDPOINT_ASYNC_INIT(sim)
        Action act() override {
            return request->readBodyToDtoAsync<oatpp::Object<user_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&sim::returnResponse);
        }
        Action returnResponse(const oatpp::Object<user_dto>& body){
            OATPP_ASSERT_HTTP(body, Status::CODE_404, "error");
            temp2 = body;
            return request->readBodyToDtoAsync<oatpp::Object<user_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&gms::use2);
        }
    };
    ENDPOINT_ASYNC("POST", "/sendmessage",sms) {
    ENDPOINT_ASYNC_INIT(sms)
        Action act() override {
            return request->readBodyToDtoAsync<oatpp::Object<message_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&sms::returnResponse);
        }
        Action returnResponse(const oatpp::Object<message_dto>& body){
            OATPP_ASSERT_HTTP(body, Status::CODE_404, "error");
            me = body;
        }
    };
    ENDPOINT_ASYNC("POST", "/add1",af1) {
    ENDPOINT_ASYNC_INIT(af1)
        Action act() override {
            return request->readBodyToDtoAsync<oatpp::Object<user_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&af1::use);
        }
        Action use(const oatpp::Object<user_dto>& body){
            OATPP_ASSERT_HTTP(body, Status::CODE_404, "error");
             temp3 = body;
        }
    };
    ENDPOINT_ASYNC("POST", "/add2",af2) {
    ENDPOINT_ASYNC_INIT(af2)
        Action act() override {
            return request->readBodyToDtoAsync<oatpp::Object<user_dto>>(
                    controller->getDefaultObjectMapper()
            ).callbackTo(&af2::returnResponse);
        }
        Action returnResponse(const oatpp::Object<user_dto>& body){
            OATPP_ASSERT_HTTP(body, Status::CODE_404, "error");
            OATPP_ASSERT_HTTP(temp3, Status::CODE_404, "error");
            temp4 = body;
            /*
             *
             * add friend
             * code ......
             *
             * */
            return _return(controller->createDtoResponse(Status::CODE_100,me));
        }
    };
};
#include OATPP_CODEGEN_END(ApiController)
#endif //WEBBACK_CONTROLLER_HPP
