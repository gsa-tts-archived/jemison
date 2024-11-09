local E = import 'extract.libsonnet';
local F = import 'fetch.libsonnet';
local P = import 'pack.libsonnet';
local S = import 'serve.libsonnet';
local W = import 'walk.libsonnet';

local VCAP = import 'vcap_services.libsonnet';

{
  // :: means "not visible in the output"
  EIGHT_SERVICES: {
    extract: E.container,
    fetch: F.container,
    pack: P.container,
    serve: S.container,
    walk: W.container,
  },
  VCAP_SERVICES: VCAP.VCAP_SERVICES,
}
