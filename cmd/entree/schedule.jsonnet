// # https://crontab.guru/#0_0_*_*_0
// local W = "0 0 * * 0";
// # https://crontab.guru/#0_0_1,15_*_*
// local BW = "0 0 2,16 * *";
// # https://crontab.guru/#0_0_1_*_*
// local M = "0 0 1 * *";
// # https://crontab.guru/#0_0_1_1,4,7,10_*
// local Q = "0 0 1 1,4,7,10 *";

local wg(host, path, scheme='https') = {
  scheme: scheme,
  host: host + '.gov',
  path: path,
};

local root = '/';

{
  minutely: [
  ],
  hourly: [
    wg('coldcaserecords', root),
    wg('www.fac', root),
  ],
  daily: [
    wg('digitalcorps.gsa', root),
  ],
  weekly: [
    wg('fac', root),
    wg('cloud', root),
  ],
  monthly: [
    wg('search', root),
  ],
  quarterly: [
    wg('gsa', root),
  ],
}
