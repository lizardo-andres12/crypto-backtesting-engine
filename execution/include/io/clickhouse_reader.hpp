#pragma once

#include <clickhouse/client.h>
#include <memory>

#include "types.hpp"

class ClickhouseReader {
private:
    std::shared_ptr<clickhouse::Client> m_client_ptr;

public:
    ClickhouseReader(std::shared_ptr<clickhouse::Client> client_ptr) : m_client_ptr(client_ptr) {}

    void stream_data(
        const std::string& symbol,
        uint64_t start_ts,
        uint64_t end_ts,
	std::function<void(const std::vector<Candle>&)> callback
    ) const;
};
