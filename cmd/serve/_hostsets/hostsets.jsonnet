local gov = import 'gov/gov.jsonnet';

assert gov.number_of_sets == std.length(gov.hostsets);

// Looks at a single hostset.
// This is an array [{}, {}, ...]
local validate(set) =
  [
    [
      std.objectHas(o, f)
      for f in ['scheme', 'host', 'paths']
    ],
    for o in set
  ];

{
  validated: [
    validate(subset)
    for subset in [set for set in gov.hostsets]
  ],
  gov: gov.hostsets,
}
