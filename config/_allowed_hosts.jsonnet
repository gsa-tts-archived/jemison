local util = import 'domain64/util.libsonnet';

local lookup = {
  all: [
    0,
    18446744073709551616,
  ],
  three: [
    [
      72057748656750592,
      72057748656750592,
    ],
    [
      72058358542106624,
      72058358542106624,
    ],
    [
      72057911865508096,
      72057911865508096,
    ],
    [
      72058358542106624,
      72058362837073664,
    ],
  ],
  nih: [
    util.toDec('0100008D00000000'),
    util.toDec('0100008DFFFFFF00'),
  ],
  uscg: [
    util.toDec('0300002000000000'),
    util.toDec('0300002FF0000000'),
  ],
  spaceforce: [
    util.toDec('0300001E00000000'),
    util.toDec('0300001EFF000000'),
  ],
  nasa: [
   util.toDec("0100008700000000"),
   util.toDec("01000087FF000000")
  ],
  eighteeneff: [
   util.toDec("0100000100000000"),
   util.toDec("01000001F0000000")
  ],
  digital: [
   util.toDec("0100003700000000"),
   util.toDec("01000037F0000000")
  ],
  fedramp: [
   util.toDec("0100005000000000"),
   util.toDec("0100005F00000000")
  ],
  dec15: [
    self.nih 
    + self.three 
    + self.uscg 
    + self.spaceforce 
    + self.nasa
    + self.eighteeneff 
    + self.digital 
    + self.fedramp
  ],
};

local toS(o) = if 
  std.isArray(o) then "[" + std.foldl(function(prev_result, e) prev_result + " " + toS(e), o, "") + " ]"
  else std.toString(o);

assert std.isArray(std.objectFields(lookup));

local debug(key, o) = std.trace(key + ": " + toS(o), o);

{
  [k]: { name: k, low: lookup[k][0], high: lookup[k][1] } 
  for k
  in debug("fields", std.objectFields(lookup))
}