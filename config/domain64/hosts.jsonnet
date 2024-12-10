local d64 = import 'domain64.jsonnet';

// Used to generate a starting schedules.libsonnet.
// Probably should not be used again.

{
  [rfqdn]: 'monthly'
  for rfqdn
  in
    std.flattenArrays(
      [d64[key].Domain64s for key in std.objectFields(d64)]
    )
}
