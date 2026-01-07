#pragma once

#include <concepts>
#include <string>

#include "types.hpp"

// Compile-time validator for trading strategies.
// All requirements MUST be satisfied to ensure compatibility with the execution engine.
template <typename T>
concept IsStrategy = requires(T strategy, const Candle& candle) {
    // Requirement: Must have an 'on_candle' method that takes a const Candle reference.
    // This is where you should put your main strategy logic as this will run once per candle.
    { strategy.on_candle(candle) } -> std::same_as<Signal>;

    // Requirement: Must have a 'get_name' for logging purposes.
    { strategy.get_name() } -> std::convertible_to<std::string>;
};

