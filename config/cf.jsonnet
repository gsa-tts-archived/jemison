local B = import 'base.jsonnet';

local E = import 'extract.jsonnet';
local F = import 'fetch.jsonnet';
local P = import 'pack.jsonnet';
local S = import 'serve.jsonnet';

{
  // :: means "not visible in the output"
  EIGHT_SERVICES: {
    extract: E.cf,
    fetch: F.cf,
    pack: P.cf,
    serve: S.cf,
    }
}
