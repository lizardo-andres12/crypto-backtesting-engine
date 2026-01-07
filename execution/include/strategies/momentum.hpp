#pragma once

#include <cstdint>
#include <deque>
#include <iostream>
#include <string>

#include "types.hpp"


template <size_t WindowSize, uint64_t Capital>
class MomentumStrategy {
private:
    std::deque<uint64_t> m_price_window;
    uint64_t m_current_capital = Capital;
    uint64_t m_current_sum = 0;

public:
    std::string get_name() const {
        return "Momentum_v1";
    }

    Signal on_candle(const Candle& candle) {
        // Update our window and sum for the moving average
        m_price_window.push_back(candle.close);
        m_current_sum += candle.close;

        // Wait for window to fill with ticks
        if (m_price_window.size() < WindowSize) return;

        if (m_price_window.size() > WindowSize) {
            uint64_t removed = m_price_window.front();
            m_price_window.pop_front();

            m_current_sum -= removed;
        }

        // Calculate simple moving average (lossy due to integer divison)
        uint64_t sma = m_current_sum / WindowSize;

        // Write results
        if (candle.close > sma) {
            return Signal{
                .timestamp = candle.timestamp,
                .action = Action::BUY,
                .price = candle.close,
                .size = 1.0,
            };
        } else {
            std::cout << "SELL @ " << candle.close << " SMA: " << sma << '\n';
        }
    }
};
