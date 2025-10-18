# transit-watcher

An attempt to sync [https://github.com/catouberos/transit-radar](transit-radar) with data from GoBus and MultiGo

## Workflows

### Importer

## Architecture Decisions

### Event Consumption

As this service fetch data periodically from two different service providers 
and they have one-way relationship in term of data dependencies (e.g. vehicle
real-time data depends on its route-variant and vehicle ID), event consumption
does not gurantee to succeed due to multiple reasons:

1. Upon first-run, required data might not be initialized (e.g. real-time position
data might be polling, but cannot be consume due to missing route-variant relationships)
2. One provider might contain outdated data, which is not reflected to be used by
another provider

Therefore, consumption is delayed and retried until success with a specific
retry policy.
