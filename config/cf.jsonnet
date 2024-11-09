local E = import 'extract.libsonnet';
local F = import 'fetch.libsonnet';
local P = import 'pack.libsonnet';
local S = import 'serve.libsonnet';

{
  // :: means "not visible in the output"
  EIGHT_SERVICES: {
    extract: E.cf,
    fetch: F.cf,
    pack: P.cf,
    serve: S.cf,
    }
}
