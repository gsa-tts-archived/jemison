local B = import 'base.jsonnet';
local service = 'extract';

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
    'extract_pdf',
    { cf: true, container: true },
  ],
  [
    'extract_html',
    { cf: true, container: true },
  ],
  [
    'walkabout',
    { cf: true, container: true },
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
