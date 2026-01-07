#pragma once

#include <cstdint>


class SimpleAggregator {
private:
    uint64_t total_pnl;
    double sharpe_ratio;
    double max_drawdown;
};
