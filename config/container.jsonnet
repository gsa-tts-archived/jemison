local E = import 'extract.jsonnet';
local F = import 'fetch.jsonnet';
local P = import 'pack.jsonnet';
local S = import 'serve.jsonnet';
local W = import 'walk.jsonnet';

local VCAP = import 'vcap_services.jsonnet';

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
