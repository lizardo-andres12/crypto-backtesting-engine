# IMPORTANT

Engine flow:
- Control plane writes the static configuration file `config.hpp` to the include directory (symbol, start/end times, starting capital, strategy, ...) and the main file of the execution unit is compiled
- The engine requests the data as a stream from clickhouse
- On a single candle:
    - The engine calls the strategy to process the candle and generate a new signal.
    - The signal, containing the corresponding action and amount, is returned for an aggregator to update the metrics of the strategy.
- After processing all the data, the engine is responsible for writing the final results back to a local db and the signal history (interval play-by-play of the strategy) to a local file for replay by a visualization service or copying to artifact storage. **NOTE** These are sliced results, and for a full backtest overview, all files should be concatenated, and the final