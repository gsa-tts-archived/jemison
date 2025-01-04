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
    { cf: 'warn', container: 'debug', localhost: 'debug' },
  ],
] + B.parameters;

{
  creds:: [[service] + x for x in credentials],
  params:: [[service] + x for x in parameters],
  cf: { name: service } +
      B.params('credentials', 'cf', service, self.creds) +
      B.params('parameters', 'cf', service, self.params),
  container: { name: service } +
             B.params('credentials', 'container', service, self.creds) +
             B.params('parameters', 'container', service, self.params),
}
