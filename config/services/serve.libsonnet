local B = import 'base.libsonnet';
local service = 'serve';

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
    'debug_level',
    { cf: 'warn', container: 'warn' },
  ],
  [
    'external_port',
    { cf: 443, container: 10000 },
  ],
  [
    'external_scheme',
    { cf: 'https', container: 'http' },
  ],
  [
    'external_host',
    { cf: 'jemison.app.cloud.gov', container: 'localhost' },
  ],
  [
    'static_files_path',
    {
      cf: '/home/vcap/app/static',
      container: '/home/vcap/app/assets/static',
    },
  ],
  [
    'database_files_path',
    {
      cf: '/home/vcap/app/assets/databases',
      container: '/home/vcap/app/assets/databases',
    },
  ],
  [
    'results_per_query',
    { cf: 10, container: 10 },
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
