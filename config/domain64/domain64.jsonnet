local com = import 'com.libsonnet';
local edu = import 'edu.libsonnet';
local gov = import 'gov.libsonnet';
local mil = import 'mil.libsonnet';
local net = import 'net.libsonnet';
local org = import 'org.libsonnet';
local schedules = import 'schedules.libsonnet';

local assertion = import 'assertions.libsonnet';
local util = import 'util.libsonnet';

local allFQDN(tld, domains) = std.flattenDeepArray([
  [
    (if kid_val != '' then kid_val + '.' else '') + std.get(domains, d).name + '.' + tld
    for kid_val in std.objectValues(std.get(domains, d).children)
  ]
  for d in std.objectFields(domains)
]);

local allRFQDN(tld, domains) = std.flattenArrays([
  [
    std.filter(std.isString, [tld, std.get(domains, d).name, (if kid_val != '' then kid_val else null)])
    for kid_val in std.objectValues(std.get(domains, d).children)
  ]
  for d in std.objectFields(domains)
]);

local allDomain64(tld, domains) = std.flattenDeepArray([
  [
    // gov
    util.getTLD(tld) + d + kid_key + '00'
    for kid_key in std.objectFields(std.get(domains, d).children)
  ]
  for d in std.objectFields(domains)
]);

local fqdnToDomain64(allFQDN, allDomain64) = {
  [pair[0]]: pair[1]
  for pair in std.mapWithIndex(function(ndx, d64) [allFQDN[ndx], d64], allDomain64)
};

local domain64ToFqdn(allFQDN, allDomain64) = {
  [pair[1]]: pair[0]
  for pair in std.mapWithIndex(function(ndx, d64) [allFQDN[ndx], d64], allDomain64)
};

local fqdnToSchedule(tld, fqdns, schedules) = {
  [pair[0]]: pair[1]
  for pair
  in
    std.filter(function(p) std.endsWith(p[0], tld),
               std.flattenArrays(
                 [
                   [[fqdn, key] for fqdn in schedules[key]]
                   for key in std.objectFields(schedules)
                 ]
               ))
};

local setDomainsKeyF(pair) = pair[0];

local domains(fqdns, d64s) =
  {
    [pair[0]]: pair[1]
    for pair
    in std.set(std.mapWithIndex(function(ndx, d64)
      [
        util.fqdnTLD(fqdns[ndx]) + '.' + util.fqdnDomain(fqdns[ndx]),
        std.slice(d64, 0, 8, 1),
      ], d64s), keyF=setDomainsKeyF)
  };

// std.set([std.slice(x, 0, 8, 1) for x in d64s]);

local tld_arr = ['gov', 'mil', 'com', 'net', 'edu', 'org'];
local domain_arr = [gov.domains, mil.domains, com.domains, net.domains, edu.domains, org.domains];

{
  [pair[0]]: {
    FQDNs: allFQDN(pair[0], pair[1]),
    RFQDNs: std.map(function(ls) std.join('.', ls), allRFQDN(pair[0], pair[1])),
    Domain64s: allDomain64(pair[0], pair[1]),
    FQDNToDomain64: fqdnToDomain64(self.FQDNs, self.Domain64s),
    Domain64ToFQDN: domain64ToFqdn(self.FQDNs, self.Domain64s),
    RDomainToDomain64: domains(self.FQDNs, self.Domain64s),
    Schedules: {
      [d64]: schedules[d64]
      for d64 in self.Domain64s
    },
    assert assertion.andMap([std.length(d64) == 16 for d64 in allDomain64(pair[0], pair[1])]),
    assert assertion.andMap([assertion.validateDomains(domains) for domains in domain_arr]),
  }
  for
  pair
  in std.mapWithIndex(function(ndx, d) [tld_arr[ndx], d], domain_arr)
} + {
  TLDs: tld_arr,
}
