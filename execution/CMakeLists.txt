cmake_minimum_required(VERSION 3.20)
project(execution_engine LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED True)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

add_compile_options(-march=native -Wall)

include(FetchContent)

FetchContent_Declare(
    clickhouse-cpp
    GIT_REPOSITORY https://github.com/ClickHouse/clickhouse-cpp.git
    GIT_TAG v2.6.0
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/lib/clickhouse_cpp
    DOWNLOAD_EXTRACT_TIMESTAMP FALSE
)

FetchContent_MakeAvailable(clickhouse-cpp)

add_executable(execution_engine
    src/main.cpp 
    src/io/clickhouse_reader.cpp
)

target_include_directories(execution_engine PRIVATE include)
target_include_directories(execution_engine PRIVATE lib)

target_link_libraries(execution_engine PRIVATE clickhouse-cpp-lib)

