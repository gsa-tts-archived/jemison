// For performance, we'll pull the rendered JSON file.
// This means this file wants to be built *after* the domain64 file.
// domain64.jsonnet is expensive to execute.
local d64 = import '../domain64/domain64.json';
local util = import '../domain64/util.libsonnet';

local fqdn_include = [
  'cloud.gov',
  'search.gov',
  'www.fac.gov',
];

local ranges_to_include = [
  ['010000B200000000', '010000B2FFFFFF00'],
];

local getTLD(fqdn) = std.reverse(std.split(fqdn, '.'))[0];

assert getTLD('cloud.gov') == 'gov';

local lookupD64(fqdn) = d64[getTLD(fqdn)].FQDNToDomain64[fqdn];

local d64ToDec(d) = std.parseHex(d);

assert lookupD64('cloud.gov') == '0100002400000000';

{
  three: [[d64ToDec(lookupD64(fqdn)), d64ToDec(lookupD64(fqdn))] for fqdn in fqdn_include] +
         [[d64ToDec(p[0]), d64ToDec(p[1])] for p in ranges_to_include],
  all: [[0, d64ToDec('FFFFFFFFFFFFFF00')]],
  nih: [
    [
      util.toDec('0100008D00000000'),
      util.toDec('0100008DFFFFFF00'),
    ],
  ],

  uscg: [
    [
      util.toDec('0200002000000000'),
      util.toDec('0200002FF0000000'),
    ],
  ],
  spaceforce: [
    [
      util.toDec('0200001E00000000'),
      util.toDec('0200001EFF000000'),
    ],
  ],
  nasa: [
    [
      util.toDec('0100008700000000'),
      util.toDec('01000087FF000000'),
    ],
  ],
  eighteeneff: [
    [
      util.toDec('0100000100000000'),
      util.toDec('01000001F0000000'),
    ],
  ],
  digital: [
    [
      util.toDec('0100003700000000'),
      util.toDec('01000037F0000000'),
    ],
  ],
  fedramp: [
    [
      util.toDec('0100005000000000'),
      util.toDec('0100005F00000000'),
    ],
  ],
  dec15: self.nih
         + self.three
         + self.nasa
         + self.uscg
         + self.spaceforce
         + self.eighteeneff
         + self.digital
         + self.fedramp,
}
