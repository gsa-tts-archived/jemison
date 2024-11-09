local E = import 'extract.libsonnet';
local F = import 'fetch.libsonnet';
local P = import 'pack.libsonnet';
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
      E.container,
      F.container,
      P.container,
      S.container,
      V.container,
      W.container,
    ],
  },
  VCAP_SERVICES: VCAP.VCAP_SERVICES,
}
