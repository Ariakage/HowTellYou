//
// Created by bishijie on 23-7-10.
//

#ifndef WEBBACK_DBCOMPONENT_HPP
#define WEBBACK_DBCOMPONENT_HPP
#include "DbClient.cpp"
#include "oatpp-sqlite/orm.hpp"
class dbcomponent {
public:
    OATPP_CREATE_COMPONENT(std::shared_ptr<user_client>, myDatabaseClient)([] {

        /* Create database-specific ConnectionProvider */
        auto connectionProvider = std::make_shared<oatpp::sqlite::ConnectionProvider>("/path/to/database.sqlite");

        /* Create database-specific Executor */
        auto executor = std::make_shared<oatpp::sqlite::Executor>(connectionProvider);

        /* Create MyClient database client */
        return std::make_shared<MyClient>(executor);

    }());

};
#endif //WEBBACK_DBCOMPONENT_HPP
