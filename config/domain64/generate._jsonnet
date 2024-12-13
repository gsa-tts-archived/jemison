local gov = import 'domains.libsonnet';
local hosts = import '../schedules/hosts.libsonnet';
local util = import 'util.libsonnet';

local justHosts = std.flattenArrays([hs for hs in std.objectValues(hosts)]);

/*
local domains = {
  '000001': {
    name: 'search',
    children: {
      '000000': "",
    },
  },
  */

local getKids(rfqdn, tld, d) = 
  [std.join(".", sub[1:]) 
  for sub in std.filter(function(o) o[0] == d, 
  [d[1:] for d in std.filter(function(o) o[0] == tld, rfqdn)])];

local makeChildren(kids) = {[pair[0]]: pair[1] for pair in 
  std.filter(function(p) p[1] != "", std.mapWithIndex(function(ndx, o) [util.toHex(ndx+1), o], kids))
};

local byTLD(tld) = {
  jh:: justHosts,
  rfqdn:: [std.reverse(std.split(h, ".")) for h in $.jh],
  justGov:: std.filter(function(v) v[0] == tld, $.rfqdn),
  domains:: std.uniq(std.sort(std.uniq([d[1] for d in $.justGov]))),
  tld: {
    [pair[0]]: pair[1] 
    for pair in 
    std.mapWithIndex(function(ndx, o) [util.toHex(ndx+1), { name: o, children: makeChildren(getKids($.rfqdn, tld, o))}] , $.domains)
    },
};

byTLD("org")
