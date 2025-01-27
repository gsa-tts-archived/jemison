local A = import 'admin.libsonnet';
local EN = import 'entree.libsonnet';
local EX = import 'extract.libsonnet';
local F = import 'fetch.libsonnet';
local M = import 'migrate.libsonnet';
local P = import 'pack.libsonnet';
local R = import 'resultsapi.libsonnet';
local S = import 'serve.libsonnet';
local V = import 'validate.libsonnet';
local W = import 'walk.libsonnet';

local VCAP = import 'vcap_services.libsonnet';

{
  APPENV: 'DOCKER',
  HOME: '/home/vcap/app',
  MEMORY_LIMIT: '512m',
  PWD: '/home/vcap',
  TMPDIR: '/home/vcap/tmp',
  USER: 'vcap',
  EIGHT_SERVICES: {
    'user-provided': [
      A.container,
      EN.container,
      EX.container,
      F.container,
      M.container,
      P.container,
      S.container,
      V.container,
      W.container,
      R.container,
    ],
  },
  VCAP_SERVICES: VCAP.VCAP_SERVICES(
    'minio',
    [
      ['jemison-queues-db', 'jemison-queues-db', 5432],
      ['jemison-work-db', 'jemison-work-db', 5432],
      ['jemison-search-db', 'jemison-search-db', 5432],
    ]
  ),
}
