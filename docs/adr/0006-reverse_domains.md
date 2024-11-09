# Use what's avaialble

Date: 2024-11-09

## Status

Accepted

## Context

We will encounter hosts and paths in a "forward" context, but think backwards we must.

## Decision

Hosts (domains) must always be represented in reverse domain name notation (https://en.wikipedia.org/wiki/Reverse_domain_name_notation). For example, instead of storing a domain as 

`app.fac.gov`

we should store

`gov.fac.app`

instead.

## Consequences

Extra work is required to process and flip (for example) `URL` structs in Golang to marshal/unmarshal this kind of data. However, a simple function (or struct method) along the lines of

```
func reverseDomain(domain string) string {
	parts := strings.Split(domain, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}
```

will do the job. 

This will make parent/child relationships, or the application of various rules (e.g. something matching a TLD pattern of `gov.fac.*`) much easier to implement.