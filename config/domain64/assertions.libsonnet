local and(a, b) = a && b;
local andMap(arr) = std.foldl(and, arr, true);
local isTrue(v) = v == true;

local validateDomains(domains) =
  // Assert that the keys are all length 6
  isTrue(andMap([
    std.assertEqual(std.length(v), 6)
    for v in std.objectFields(domains)
  ]))
  &&
  // And, that we have unique mappings
  std.set(std.objectFields(domains)) == std.objectFields(domains)
  &&
  std.length(std.set(std.objectFields(domains))) == std.length(std.objectFields(domains))
  &&
  // Asserting the values is redundant with the length == 6 test, but still valuable
  andMap([std.parseHex(v) <= std.parseHex('FFFFFF') for v in std.objectFields(domains)])
;

{
  and:: and,
  andMap:: andMap,
  isTrue:: isTrue,
  validateDomains:: validateDomains,
}
