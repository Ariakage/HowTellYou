#include "oatpp/web/server/HttpConnectionHandler.hpp"
#include "oatpp/network/virtual_/client/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/server/ConnectionProvider.hpp"
#include "oatpp/network/virtual_/Interface.hpp"
#include "oatpp/parser/json/mapping/ObjectMapper.hpp"
#include "oatpp/core/macro/component.hpp"
#include "controller.hpp"
#include "component.hpp"
#include "controller.hpp"
#include "component.hpp"
#include "oatpp/network/Server.hpp"
void run() {
    Component components;
    OATPP_COMPONENT(std::shared_ptr<oatpp::web::server::HttpRouter>, router);
    auto usercontroller = std::make_shared<user_controller>();
    router->addController(usercontroller);
    OATPP_COMPONENT(std::shared_ptr<oatpp::network::ConnectionHandler>, connectionHandler);
    OATPP_COMPONENT(std::shared_ptr<oatpp::network::ServerConnectionProvider>, connectionProvider);
    oatpp::network::Server server(connectionProvider, connectionHandler);
    OATPP_LOGI("WebQQ", "Server running on port %s", connectionProvider->getProperty("port").getData());
    server.run();
}

int main(int argc, const char * argv[]) {
    oatpp::base::Environment::init();
    run();
    oatpp::base::Environment::destroy();
    return 0;
}