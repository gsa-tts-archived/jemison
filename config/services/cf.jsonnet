local A = import 'admin.libsonnet';
local C = import 'collect.libsonnet';
local EN = import 'entree.libsonnet';
local EX = import 'extract.libsonnet';
local F = import 'fetch.libsonnet';
local M = import 'migrate.libsonnet';
local P = import 'pack.libsonnet';
local S = import 'serve.libsonnet';
local V = import 'validate.libsonnet';
local W = import 'walk.libsonnet';

{
  // :: means "not visible in the output"
  EIGHT_SERVICES: {
    'user-provided': [
      A.cf,
      EN.cf,
      EX.cf,
      F.cf,
      M.cf,
      P.cf,
      S.cf,
      V.cf,
      W.cf,
      C.cf,
    ],
  },
}
