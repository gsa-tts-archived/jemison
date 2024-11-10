local E = import 'extract.libsonnet';
local F = import 'fetch.libsonnet';
local P = import 'pack.libsonnet';
local S = import 'serve.libsonnet';
local V = import 'validate.libsonnet';
local W = import 'walk.libsonnet';

local VCAP = import 'vcap_services.libsonnet';

// https://github.com/GSA-TTS/jemison/blob/4ab9d0b2137384d02585a6b0b80d3685a112ea1e/config/localhost.yaml
{
  APPENV: 'LOCALHOST',
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
      S.container {
        parameters: {
          static_files_path: '../../assets/static',
          database_files_path: '../../assets/databases',
        },
      },
      V.container,
      W.container,
    ],
  },
  VCAP_SERVICES: VCAP.VCAP_SERVICES('localhost', 'localhost'),
}
