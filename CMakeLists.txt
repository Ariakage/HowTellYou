cmake_minimum_required(VERSION 3.24)
project(webback)

set(CMAKE_CXX_STANDARD 17)

add_executable(webback main.cpp controller.hpp component.hpp dto.hpp DbClient.cpp client.cpp dbcomponent.hpp)
target_link_libraries(webback
        pthread
        oatpp
        )