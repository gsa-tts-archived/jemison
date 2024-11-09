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
    { cf: 10, container: 10 },
  ],
  [
    'polite_sleep',
    { cf: 2, container: 2 },
  ],
  [
    'polite_cache_default_expiration',
    { cf: B.hours(10), container: B.minutes(10) },
  ],
  [
    'polite_cache_cleanup_interval',
    { cf: B.hours(3), container: B.minutes(5) },
  ],
];

{
  creds:: [[service] + x for x in credentials],
  params:: [[service] + x for x in parameters],
  cf: B.params('credentials', 'cf', service, self.creds) +
            B.params('parameters', 'cf', service, self.params),
  container: B.params('credentials', 'container', service, self.creds) +
                   B.params('parameters','container', service, self.params),
}
