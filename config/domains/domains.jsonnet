local a = import '../assertions.libsonnet';
local gov = import 'gov/gov.libsonnet';
local mil = import 'mil/mil.libsonnet';

local has = [
  a.andMap([std.objectHas(o, key) for o in gov])
  for key in ['id', 'tld', 'domain', 'subdomains', 'indexFrequency']
];

assert a.andMap(has);

{
  gov: gov,
  mil: mil,
}
