local B = import 'lib/builders.libsonnet';
local C = import 'lib/const.libsonnet';

local ORG = 'sandbox-gsa';
local SPACE = 'matthew.jadud';
local S3_PRIVATE_BUCKETS = ['extract', 'fetch', 'serve'];


local _id(service, type) = service + '.' + type + '_id';
local bucket_id(service) = _id('module.private-' + service, 'bucket');
local instance_id(service) = _id('module.' + service, 'instance');

local data = {
  data: {
    cloudfoundry_domain: {
      public: {
        name: 'app.cloud.gov',
      },
    },
    cloudfoundry_space: {
      app_space: {
        org_name: ORG,
        name: SPACE,
      },
    },
  },
};

local module = {
  module: B.db('queues_database', ORG, SPACE, 'jemison-queues-db', 'micro-psql') +
          B.db('work_database', ORG, SPACE, 'jemison-work-db', 'micro-psql') +
          B.s3_multi(S3_PRIVATE_BUCKETS, ORG, SPACE, 'basic'),
};

local DEFAULT_ENV = {
  ENV: 'SANDBOX',
  API_KEY: '${var.api_key}',
  DEBUG_LEVEL: 'warn',
  GIN_MODE: 'production',
};

local services = B.services([
  {
    name: 'admin',
    disk: C.DISK.M,
    ram: 256,
    instances: 1,
    bindings: [
      instance_id('queues_database'),
    ],
    routes: ['cloudfoundry_route.admin_route.id'],
    env: DEFAULT_ENV,
  },
  {
    name: 'entree',
    disk: C.DISK.M,
    ram: 512,
    instances: 1,
    bindings: [
      instance_id('work_database'),
      instance_id('queues_database'),
    ],
  },
  {
    name: 'extract',
    disk: C.DISK.XL,
    ram: 1024,
    instances: 2,
    bindings: [
      bucket_id('fetch'),
      bucket_id('extract'),
      instance_id('queues_database'),
    ],
    env: DEFAULT_ENV,
  },
  {
    name: 'fetch',
    disk: C.DISK.L,
    ram: 1024,
    instances: 1,
    bindings: [
      bucket_id('fetch'),
      instance_id('work_database'),
      instance_id('queues_database'),
    ],
    env: DEFAULT_ENV,
  },
  {
    name: 'pack',
    disk: C.DISK.XL,
    ram: 512,
    instances: 1,
    bindings: [
      bucket_id('extract'),
      bucket_id('serve'),
      instance_id('queues_database'),
    ],
    env: DEFAULT_ENV,
  },
  {
    name: 'serve',
    disk: C.DISK.FULL,
    ram: 512,
    instances: 1,
    bindings: [
      bucket_id('serve'),
      bucket_id('fetch'),
      bucket_id('serve'),
      instance_id('queues_database'),
    ],
    routes: ['serve'],
    env: DEFAULT_ENV,
  },
  {
    name: 'walk',
    disk: C.DISK.L,
    ram: 1024,
    instances: 2,
    bindings: [
      bucket_id('fetch'),
      instance_id('queues_database'),
    ],
    env: DEFAULT_ENV,
  },
]);

local resource = {
  resource: services
            {
    cloudfoundry_route: {
      admin_route: {
        space: 'data.cloudfoundry_space.app_space.id',
        domain: 'data.cloudfoundry_domain.public.id',
        hostname: 'jemison-admin',
      },
    },
  },
};

{} + data + module + resource
