local hosts = [
  'www.fac.gov',
  'cloud.gov',
  'search.gov',
];

{
  monthly: [{ scheme: 'https', path: '/', host: h } for h in hosts],
}
