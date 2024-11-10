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
    uri: 'http://' + host + ':9000',
    port: 9000,
    insecure_skip_verify: false,
    access_key_id: 'numbernine',
    secret_access_key: 'numbernine',
    region: 'us-east-1',
    bucket: service,
    endpoint: host + ':9000',
    fips_endpoint: 'http://' + host + ':9000',
    additional_buckets: [],
  },
  syslog_drain_url: 'https://ALPHA.drain.url',
  volume_mounts: ['no_mounts'],
};

local rds(db, host) = {
  label: db,
  provider: null,
  plan: null,
  name: db,
  tags: ['aws', 'rds', 'postgres'],
  instance_guid: std.substr(std.base64(self.label+host), 0, 16),
  binding_guid: std.substr(std.base64(self.credentials.uri), 0, 16),
  binding_name: null,
  instance_name: db,
  credentials: {
    db_name: 'postgres',
    host: host,
    name: 'postgres',
    password: '',
    port: 5432,
    username: 'postgres',
    uri: 'postgresql://' + self.username + '@' + host + ':' + self.port + '/' + self.db_name + '?sslmode=disable',
  },
  syslog_drain_url: 'https://BRAVO.drain.url',
  volume_mounts: ['no_mounts'],
};

local VCAP_SERVICES(host, db_host) = {
    s3: [
      s3('extract', host),
      s3('fetch', host),
      s3('serve', host),
    ],
    'aws-rds': [rds('jemison-db', host)]
};

{
  VCAP_SERVICES:: VCAP_SERVICES,
}
