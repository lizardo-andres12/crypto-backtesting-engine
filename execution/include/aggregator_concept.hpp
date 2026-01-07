#pragma once

#include <concepts>
#include <string>

#include "types.hpp"

template <typename T>
concept IsAggregator = requires(T aggregator, const Signal& signal) {
    // Requirement: Must have an 'on_signal' method that takes a const Signal reference and
    // updates strategy performance metrics accordingly.
    { aggregator.on_signal(signal) } -> std::same_as<void>;

    // Requirement: Must have a 'output_stats' method to write final stats.
    // All defined stats are fields in BacktestResults struct.
    { aggregator.output_metrics() } -> std::same_as<BacktestResult>;

    // Requirement: Must have a 'get_name' for logging purposes.
    { aggregator.get_name() } -> std::convertible_to<std::string>;
};

