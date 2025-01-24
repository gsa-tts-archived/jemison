local s3(service, host) = {
  label: 's3-' + service,
  provider: 'minio-local',
  plan: 'basic',
  name: service,
  tags: ['aws', 's3', 'object-storage'],
  instance_guid: std.substr(std.base64(self.label+service), 0, 16),
  binding_guid: std.substr(std.base64(self.provider+service), 0, 16),
  binding_name: null,
  instance_name: service + '-storage',
  credentials: {
    uri: 'http://' + host + ':9100',
    port: 9100,
    insecure_skip_verify: false,
    access_key_id: 'numbernine',
    secret_access_key: 'numbernine',
    region: 'us-east-1',
    bucket: service,
    endpoint: host + ':9100',
    fips_endpoint: 'http://' + host + ':9100',
    additional_buckets: [],
  },
  syslog_drain_url: 'https://ALPHA.drain.url',
  volume_mounts: ['no_mounts'],
};

local rds(db_host, db_name, port) = {
  label: db_name,
  provider: null,
  plan: null,
  name: db_name,
  tags: ['aws', 'rds', 'postgres'],
  instance_guid: std.substr(std.base64(self.label+db_name), 0, 16),
  binding_guid: std.substr(std.base64(self.credentials.uri), 0, 16),
  binding_name: null,
  instance_name: db_name,
  credentials: {
    db_name: 'postgres',
    host: db_host,
    name: 'postgres',
    password: '',
    port: port,
    username: 'postgres',
    uri: 'postgresql://' + self.username + '@' + db_host + ':' + self.port + '/' + self.db_name + '?sslmode=disable',
  },
  syslog_drain_url: 'https://BRAVO.drain.url',
  volume_mounts: ['no_mounts'],
};

local VCAP_SERVICES(s3_host, db_hosts_and_names) = {
    s3: [
      s3('extract', s3_host),
      s3('fetch', s3_host),
      s3('serve', s3_host),
    ],
    'aws-rds': [rds(hnp[0], hnp[1], hnp[2],) for hnp in db_hosts_and_names],
};

{
  VCAP_SERVICES:: VCAP_SERVICES,
}
