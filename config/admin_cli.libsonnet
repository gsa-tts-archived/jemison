local B = import 'base.libsonnet';
local service = 'admin';

local credentials = [
  [
    'port',
    { cf: 8080, container: 8888, localhost: 8888 },
  ],
];

local parameters = [
  [
    'debug_level',
    { cf: 'warn', container: 'debug', localhost: 'debug'},
  ],
  [
    'api_base_url',
    {
      cf: 'NA',
      container: 'http://admin:8888',
      localhost: 'http://localhost:10001/api',
    },
  ],
];

{
  creds:: [[service] + x for x in credentials],
  params:: [[service] + x for x in parameters],
  cf: B.params('credentials', 'cf', service, self.creds) +
      B.params('parameters', 'cf', service, self.params),
  localhost: { name: service } +
             B.params('credentials', 'localhost', service, self.creds) +
             B.params('parameters', 'localhost', service, self.params),
}
