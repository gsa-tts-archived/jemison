# Use what's avaialble

Date: 2024-11-09

## Status

Accepted

## Context

Configuration is needed everywhere. This sometimes come up as YAML, sometimes JSON. 

It turns out, YAML is considered a superset of JSON, meaning all JSON documents are valid YAML. (Or, they should be.)


## Decision

[Jsonnet](https://jsonnet.org) is a small programming language that, for all intents and purposes, does nothing but emit JSON. It has mechanisms for eliminating duplication, and (more importantly) compiles Jsonnet documents into JSON, meaning that a layer of checking is applied. 

All configuration that leads to YAML or JSON will be expressed as Jsonnet (if possible at build-time), adding a layer of linting, formatting, and checking to our build process for all configuration documents.

## Consequences

A small amount of work is needed to learn Jsonnet. However, it is (for all intents and purposes) a cleaner JSON. The curve is short and shallow, and the benefits are many.