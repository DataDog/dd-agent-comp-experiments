# Trace Agent

(temporary notes on how trace-agent is structured)

Creates all the other stuff.

Has parallel workers that
 - read from receivers
 - filter against Blacklister
 - replace with Replacer
 - sample with *Sampler and EventProcessor?
 - send to TraceWriter
 - send stats to Concentrator

# Receiver

pkg/trace/api.HTTPReceiver

Listens for spans, sends them to a chan *Payload

Also sends stats to ClientStatsAggregator

# OTLPReceiver

pkg/trace/api.OTLPReceiver

Similar, but otel collector

# Concentrator

Sends stats to StatsWriter via chan

# ClientStatsAggregator

Sends stats to StatsWriter via chan

# Blacklister

Pretty simple

# Replacer
# PrioritySampler
# ErrorsSampler
# RareSampler
# NoPrioritySampler
# EventProcessor
# TraceWriter
# StatsWriter

Writes stats


---

## Dep graph

pipeline dependencies should be downstream (source -> a -> b -> c -> destination)

Agent: HTTPReceiver, OTLPReceiver, Processor, TraceWriter, StatsWriter
  - is this even necessary?

Processor: Blacklister, Replacer, Sampler, EventProcessor, TraceWriter, StatsWriter
  - multiple workers

Blacklister, etc:
  - threadsafe methods
