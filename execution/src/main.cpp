#include <cstdlib>
#include <iostream>
#include <memory>

#include "clickhouse/client.h"
#include "engine.hpp"
#include "io/clickhouse_reader.hpp"
#include "strategies/momentum.hpp"


static constexpr std::string DEFAULT_HOST = "localhost";
static constexpr uint16_t DEFAULT_PORT = 9000;
static constexpr std::string DEFAULT_USER = "default";
static constexpr std::string DEFAULT_PASSWORD = "default";
static constexpr std::string DEFAULT_DATABASE = "crypto";


// Construct options struct from static config
clickhouse::ClientOptions getOptions() {
    clickhouse::ClientOptions opts;
    opts.SetHost(DEFAULT_HOST);
    opts.SetPort(DEFAULT_PORT);
    opts.SetUser(DEFAULT_USER);
    opts.SetPassword(DEFAULT_PASSWORD);
    opts.SetDefaultDatabase(DEFAULT_DATABASE);
    
    return opts;
}

int main(int argc, char** argv) {
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    if (argc < 4) {
	    std::cerr << "Usage: ./execution_engine [symbol] [start_time (unix epoch)] [end_time (unix epoch)]\n";
	    return 1;
    }

    const std::string symbol{argv[1]};
    const uint64_t start_ts = std::strtoull(argv[2], nullptr, 10);
    const uint64_t end_ts = std::strtoull(argv[3], nullptr, 10);

    const std::shared_ptr<clickhouse::Client> client_ptr = std::make_shared<clickhouse::Client>(getOptions());
    const ClickhouseReader reader{client_ptr};
    
    BacktestEngine<MomentumStrategy<5>> engine;

    std::cout << "Starting backtest\n";
    reader.stream_data(
        symbol,
	    start_ts,
        end_ts,
	    [&](const std::vector<Candle>& batch) -> void {
	        engine.run(batch);
	    }
    );
    std::cout << "Backtest complete\n";

    return 0;
}
