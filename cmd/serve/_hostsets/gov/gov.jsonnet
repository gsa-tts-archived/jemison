local B = import 'base.libsonnet';
local cloud = import 'cloud.libsonnet';
local fac = import 'fac.libsonnet';
local search = import 'search.libsonnet';

local hostsets = [
  [B.simple('fac.gov'), search],
  [B.simple('cloud.gov')],
];

{
  hostsets: hostsets,
  number_of_sets: std.length(hostsets),
}
