local assertion = import 'assertions.libsonnet';

local lookup = ['A', 'B', 'C', 'D', 'E', 'F'];
//local toHex(v) = if v > 15 then toHex(std.floor(v / 16)) + toHex(v % 16) else if v < 10 then '%d' % v else lookup[v - 10];

local toHex(v, len=6) = std.format('%' + '0' + len + 'X', v);

local debug(obj) = std.trace(std.toString(obj), obj);

assert toHex(15) == '00000F';
assert toHex(255, 2) == 'FF';
assert toHex(12345) == '003039';

local getDomain(domain, domains) = std.filter(
  function(v) v != null,
  [if std.get(std.get(domains, key), 'name') == domain then key else null for key in std.objectFields(domains)]
)[0];


local toDec(hex) = 
  std.foldl(function(v, prev) v + prev, 
  std.mapWithIndex(function(ndx, cp) (if cp < 65 then cp-48 else cp-55)*std.pow(16, ndx), 
  std.reverse(std.map(function(s) std.codepoint(s), std.stringChars(hex)))),
  0);

assert toDec("A000FF") == 10486015;
assert toDec("0100008D00000100") == 72058199628316928;

local tlds = {
  gov: '01',
  mil: '02',
  com: '03',
  edu: '04',
  net: '05',
  org: '06',
};

// Assert that the keys are all length 2
assert assertion.isTrue(assertion.andMap([
  std.assertEqual(std.length(v), 2)
  for v in std.objectValues(tlds)
]));

// And, that we have unique mappings
assert std.set(std.objectFields(tlds)) == std.objectFields(tlds);
assert std.length(std.set(std.objectValues(tlds))) == std.length(std.objectValues(tlds));

local getTLD(tld) = std.get(tlds, tld, default=null);

assert getTLD('gov') == '01';

local fqdnTLD(s) = std.reverse(std.split(s, "."))[0];
local fqdnDomain(s) = std.reverse(std.split(s, "."))[1];

{
  getTLD:: getTLD,
  toHex:: toHex,
  debug:: debug,
  getDomain:: getDomain,
  fqdnTLD:: fqdnTLD,
  fqdnDomain:: fqdnDomain,
  toDec:: toDec,
  a1: std.trace("toDec", toDec("A000FF")),

}
