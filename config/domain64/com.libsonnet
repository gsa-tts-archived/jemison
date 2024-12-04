local assertions = import 'assertions.libsonnet';

local domains = {
  '000001': {
    children: {
      '000001': 'treasury-testing',
      '000002': 'staging.arc',
    },
    name: 'alextom',
  },
  '000002': {
    children: {
      '000001': 'my',
      '000002': 'www',
    },
    name: 'goarmy',
  },
};

assert assertions.validateDomains(domains);

{
  domains: domains,
}
