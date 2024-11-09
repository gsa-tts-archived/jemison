local B = import 'base.libsonnet';
local service = 'pack';

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
    'packing_timeout_seconds',
    { cf: B.minutes(10), container: B.minutes(3) },
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
