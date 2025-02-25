# TITLE

Date: DATE

## Status

Accepted

## Context

Jemison will initially deploy driven entirely by configuration files.

The existing design, in JSonnet, organically emerged. However, it looks neither intentional nor maintainable/scalable.

## Decision

A table-driven design, that can be edited and verified via tooling, makes more sense. The configuration should be transformed from something people can work with into something that the application can work with. That implies 

CSV -> ... -> ... -> JSON

What follows is a proposed design, for discussion, to take into account the needs and challenges we are aware of.

### The problem

Jemison needs configuration for multiple reasons.

1. **Services**. Jemison is composed of multiple services, each of which may have any number of tunable parameters. 
2. **Environments**. Jemison runs in multiple environments (local, test (e.g. Circle/CI), `dev`, `staging`, and `prod`). There may be different configuration values for each of those environments.
3. **Domains**. Jemison crawls content on domains. We need to be able to identify the domains that we expect Jemison to crawl.
4. **Domain configuration**. Any given domain might want have parameters or filters that need to be set. For example, how often do we crawl a given domain? Are there subdomains we should ignore (but there is no `sitemap.xml` to guide us)? In some respects, the domain configuration is a critical part of the crawler configuration.

At the beginning, we need to be able to edit and maintain these configurations. We want them to be 

1. Robust.
2. Verifiable/testable.
3. Computable.

But, we also want them to be *human editable*. 

### Proposed solution

We'll start with JSonnet, which is where we are now. It can be easily be compiled to JSON, and can be manipulated directly, via libraries, in all the languages we are likely to work in.

## Design

Jemison effectively needs to be driven off a set of databases, which (ideally) change infrequently. We are, therefore, encoding a set of databases as JSonnet. 

Flatter data (wide tables) are easier to work with, in some ways, becuase the relationships are clear. However, they typically lead to data duplication. 

Hashes/maps/dictionaries are easier to comprehend, but lead to more linking. Those links can be verified, however.

We'll start with dictionaries.

### TLDs

Jemison needs to know about TLDs.

tlds.libsonnet:

```
{
  "gov": 1,
  "mil": 2,
}
```

### Second-level domains

We need to know about our 2LDs. We will ultimately need to know a lot about our 2LDs.

We'll do one file per 2LD. This will ultimately lead to thousands of source configuration files. However:

1. The files can be easily combined via JSonnet
2. Properties can be asserted across all of the files, as well as within files.

We will name files in a RFQDN style (Reverse Fully Qualified Domain Name).

gov.gsa.libsonnet:

```
{
  id: 1,

  tld: "gov", // a key into tlds.libsonnet
  
  // We must explicitly encode the root with the magic key `_root`. 
  subdomains: {
    _root: 0,
    www: 1,
    acquisition: 2,
  },

  // We will use Go-style camel-case.
  indexFrequency: weekly,

  // We can imagine adding per-domain 
  // configuration as needed to this design.
  ignorePaths: {
    "*":   ["/something", "*/another_thing"],
    "acquisition": ["/lots-of-pdfs", "images"],
  }

  // There might be paths we *only* crawl.
  onlyPaths: {
      "_root": ["/magic"],
      "acquisition": ["/savings", "/agile"],
  }
}
```

Assertions can then be built around these files. For example, we can assert that there is no duplication of paths in the `ingorePaths`, or in the list of subdomains, or that no two domains have the same `id` within a given `tld`. (Why? Because these are essentially namespaced values, and we need them to be unique.)

### Affiliates

Historically, we have a notion of *affiliates*. These are a way of giving a name to a group of domains. They are used when searching via the API or a SERP to indicate one or more domains to be searched.

In documentation, `site` is sometimes synonmous with `affiliate` with `domain`. We will use `affiliate` for documentation purposes.

1. An affiliate maps to one or more domains
2. An affiliate that maps to a `tld.domain` implies `tld.domain.*`. We should crawl the entire namespace of the domain and all subdomains.
3. An affiliate that maps to `tld.domain.sub` implies `tld.domain.sub.*`. We should crawl all subsubdomains as well. 
    * **NOTE**: This implies a possible change to the D64 design.
4. An affiliate may have a path. If a site has a path, we will assume that it means to crawl that domain and everything under that path. E.g. `*/path/tld.domain` *or* `*/path/tld.domain.sub`. We will not put a star under the domain or subdomain in this case. So, we do not crawl all subdomains at a given path; this is a specialization of #1 and #2.

Ultimately, affiliates are an expression of what parts of the database someone wants to search. We have a set of defined affiliates already. We will likely want to do one of several things going forward:

1. **Allowing dynamic affiliates**. We may want to allow people to define sets of domains to search as part of a query. E.g. passing a list of domains so that a search is restricted to those domains.

#### Excluded domains

There are exclusions, which we will capture in the JSonnet design above. 

Exclusions will be expressed as regular expressions.

### Generated files

We can then use this to functionally generate additional configuration within Jsonnet.

For example, we can now import these files, and generate a list of all of our domains. We can then write assertions that guarantee uniqueness of `domain64` values, and that they are all the correct size, etc.

Generated files become templated `jsonnet` files that are built via Makefile, and ultimately want library support in Jemison for loading/accessing them easily.

## Consequences

1. The configuration is now a pipelined build process. That adds steps.
2. The configuration can validate itself via JSonnet `assert` statements. This adds assurance at the end of the process.

