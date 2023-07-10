//
// Created by bishijie on 23-7-10.
//
#include "DbClient.cpp"
#include "dbcomponent.hpp"
#include <iostream>
using namespace std;
int main(){
    OATPP_COMPONENT(std::shared_ptr<user_client>, client);
    client->createUser("admin", "admin@domain.com", UserRoles::ADMIN);
    auto result = client->getUserByUsername("admin");
    auto dataset = result->fetch<oatpp::Vector<oatpp::Object<UserDto>>>();
    auto json = jsonObjectMapper.writeToString(dataset);
    std::cout << json->c_str() << std::endl;
}