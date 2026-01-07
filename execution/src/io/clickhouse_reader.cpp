#include <format>
#include <vector>

#include "clickhouse/columns/numeric.h"
#include "io/clickhouse_reader.hpp"
#include "types.hpp"

void ClickhouseReader::stream_data(const std::string& symbol, uint64_t start_ts, uint64_t end_ts,
                                   std::function<void(const std::vector<Candle>&)> callback) const {
    std::string query = std::format("SELECT "
                                    "   toUnixTimestamp64Milli(timestamp) as ts, "
                                    "   open, high, low, close, volume "
                                    "FROM market_data_1m "
                                    "WHERE symbol = '{}' "
                                    "AND timestamp >= toDateTime64({}, 3) "
                                    "AND timestamp <= toDateTime64({}, 3) "
                                    "ORDER BY timestamp ASC",
                                    symbol, start_ts / 1000.0, end_ts / 1000.0);

    m_client_ptr->Select(query, [&](const clickhouse::Block& block) {
        if (block.GetRowCount() == 0)
            return;

        auto c_time = block[0]->As<clickhouse::ColumnInt64>();
        auto c_open = block[1]->As<clickhouse::ColumnUInt64>();
        auto c_high = block[2]->As<clickhouse::ColumnUInt64>();
        auto c_low = block[3]->As<clickhouse::ColumnUInt64>();
        auto c_close = block[4]->As<clickhouse::ColumnUInt64>();
        auto c_volume = block[5]->As<clickhouse::ColumnFloat64>();

        std::vector<Candle> batch;
        batch.reserve(block.GetRowCount());

        for (size_t i{}; i < block.GetRowCount(); ++i) {
            batch.push_back(Candle{.timestamp = static_cast<uint64_t>(c_time->At(i)),
                                   .open = c_open->At(i),
                                   .high = c_high->At(i),
                                   .low = c_low->At(i),
                                   .close = c_close->At(i),
                                   .volume = c_volume->At(i)});
        }

        callback(batch);
    });
}
