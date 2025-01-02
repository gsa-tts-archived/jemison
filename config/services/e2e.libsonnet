local B = import 'base.libsonnet';
local service = 'e2e';

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
    'template_files_path',
    {
      cf: '/home/vcap/app/templates',
      container: '/home/vcap/app/cmd/e2e/templates',
    },
  ],
  [
    'debug_level',
    { cf: 'warn', container: 'debug' },
  ],
  [
    'static_files_path',
    {
      cf: '/home/vcap/app/static',
      container: '/home/vcap/app/cmd/e2e/static',
    },
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
