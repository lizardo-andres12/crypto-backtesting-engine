# IMPORTANT

Aggregators are the means of statistic calculation and storage. Their function is to receive signals generated from strategies and update aggregation
state according to those signals. The **Signal** data structure is the only source of information for the aggregators, so they must contain all relevant
fields in their schema.