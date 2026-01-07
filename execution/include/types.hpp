#pragma once

#include <cstdint>
#include <string>


// Action represents what the strategy indicates to do on new candle
enum class Action : uint8_t {
    BUY,
    SELL,
    HOLD
};

// Candle contains all data for one full interval
struct Candle {
    uint64_t timestamp;
    uint64_t open;
    uint64_t high;
    uint64_t low;
    uint64_t close;
    double volume;
};

// Signal is the output of the strategy for one interval
struct Signal {
    uint64_t timestamp;
    Action action;
    uint64_t price;
    double size;
    int64_t pnl;
};

// BacktestResults is the end highlights for one single execution unit
struct BacktestResult {
    std::string job_id;
    std::string strategy_id;
    std::string symbol;
    uint64_t start_time;
    uint64_t end_time;
    int64_t total_pnl;
    double sharpe_ratio;
    double max_drawdown;
    uint64_t execution_time_ms;
};

