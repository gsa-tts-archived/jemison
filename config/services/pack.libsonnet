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
    { cf: 10, container: 5 },
  ],
  [
    'debug_level',
    { cf: 'warn', container: 'debug' },
  ],
  [
    # This used to be really critical, because the packer was 
    # building a DB on the fly. Now, it pulls a list of objects
    # from S3, and is more of a timeout. If we build a DB
    # multiple times, it is OK. Also, we won't build more than once
    # every "timeout seconds," so really it is more of a flush.
    'packing_timeout_seconds',
    { cf: B.minutes(3), container: 15 },
  ],
] + B.parameters;

{
  creds:: [[service] + x for x in credentials],
  params:: [[service] + x for x in parameters],
  cf: B.params('credentials', 'cf', service, self.creds) +
      B.params('parameters', 'cf', service, self.params),
  container: { name: service } +
             B.params('credentials', 'container', service, self.creds) +
             B.params('parameters', 'container', service, self.params),
}
