local domains = {
  '000001': {
    children: {
      '000001': 'devoigpublicsite',
    },
    name: 'azurewebsites',
  },
  '000002': {
    children: {
      '000001': 'fasterdata',
      '000002': 'lightbytes',
      '000003': 'www',
    },
    name: 'es',
  },
  '000003': {
    children: {
      '000001': 'nasa-smd-wp',
    },
    name: 'go-vip',
  },
  '000004': {
    children: {},
    name: 'rewardsforjustice',
  },
};

{
  domains: domains,
}
