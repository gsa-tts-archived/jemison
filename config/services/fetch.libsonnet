local B = import 'base.libsonnet';
local service = 'fetch';

local credentials = [
  [
    'port',
    { cf: 8080, container: 8888 },
  ],
];

local parameters = [
  [
    'workers',
    { cf: 10, container: 50 },
  ],
  [
    'debug_level',
    { cf: 'warn', container: 'info' },
  ],
  [
    'polite_sleep',
    { cf: 2, container: 2 },
  ],
];

  // [
  //   'polite_cache_default_expiration',
  //   { cf: B.hours(10), container: B.hours(10) },
  // ],
  // [
  //   'polite_cache_cleanup_interval',
  //   { cf: B.hours(3), container: B.minutes(30) },
  // ],

{
  creds:: [[service] + x for x in credentials],
  params:: [[service] + x for x in parameters],
  cf: B.params('credentials', 'cf', service, self.creds) +
      B.params('parameters', 'cf', service, self.params),
  container: { name: service } +
             B.params('credentials', 'container', service, self.creds) +
             B.params('parameters', 'container', service, self.params),
}
