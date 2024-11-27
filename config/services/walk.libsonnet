local B = import 'base.libsonnet';
local service = 'walk';

local credentials = [
  [
    'port',
    { cf: 8080, container: 8888 },
  ],
];

local parameters = [
  [
    'workers',
    { cf: 10, container: 20 },
  ],
  [
    'debug_level',
    { cf: 'warn', container: 'warn' },
  ],
  [
    'cache-ttl',
    { cf: B.minutes(120), container: B.minutes(120) },
  ],
  [
    'polite_cache_default_expiration',
    { cf: B.hours(10), container: B.hours(10) },
  ],
  [
    'polite_cache_cleanup_interval',
    { cf: B.minutes(120), container: B.minutes(120) },
  ],
];

{
  creds:: [[service] + x for x in credentials],
  params:: [[service] + x for x in parameters],
  cf: B.params('credentials', 'cf', service, self.creds) +
      B.params('parameters', 'cf', service, self.params),
  container: { name: service } +
             B.params('credentials', 'container', service, self.creds) +
             B.params('parameters', 'container', service, self.params),
}
