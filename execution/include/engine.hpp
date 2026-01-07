#pragma once

#include <vector>

#include "types.hpp"
#include "strategy_concept.hpp"
#include "aggregator_concept.hpp"

// BacktestEngine is the hot path execution unit that processes blocks of candlestick data from clickhouse and stores interval
// output in memory as a vector of Signals.
// TODO: calculate the total blocks using the timestamps to allow for m_signals.reserve(count) optimization (This is okay since the control plane will ensure a RAM-safe partition of total work)
template <IsStrategy Strategy, IsAggregator Aggregator>
class BacktestEngine {
private:
    std::vector<Signal> m_signals;
    Strategy m_strategy;
    Aggregator m_aggregator;

public:
    BacktestEngine() = default;

    void run(const std::vector<Candle>& data_block) {
        for (const auto& candle : data_block) {
            Signal sig{m_strategy.on_candle(candle)};
	    m_aggregator.on_signal(sig);
	    m_signals.push_back(sig);
        }
    }

    // TODO: implement a finalize method that outputs a serialized Signal history file for visualization with other service.
    void finalize() const {}
};
