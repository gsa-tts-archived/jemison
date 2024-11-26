local hosts = import 'hosts.libsonnet';

{
  // This comprehension builds the schedule from the list of hosts.
  // On the outside is an object comprehension, which builds the structure
  // weekly : [ ... ]
  // for each of the schedules in the periodicity list
  [periodicity]: [{ scheme: 'https', path: '/', host: h } for h in hosts[periodicity]]
  for periodicity in ['biweekly', 'monthly', 'quarterly', 'weekly']
}
