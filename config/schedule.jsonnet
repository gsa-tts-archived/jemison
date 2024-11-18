// # https://crontab.guru/#0_0_*_*_0
// local W = "0 0 * * 0";
// # https://crontab.guru/#0_0_1,15_*_*
// local BW = "0 0 2,16 * *";
// # https://crontab.guru/#0_0_1_*_*
// local M = "0 0 1 * *";
// # https://crontab.guru/#0_0_1_1,4,7,10_*
// local Q = "0 0 1 1,4,7,10 *";

local g(host, path, scheme='https') = {
  scheme: scheme,
  host: host + '.gov',
  path: path,
};

local wg(host, path, scheme='https') = {
  scheme: scheme,
  host: 'www.' + host + '.gov',
  path: path,
};

local root = '/';

{
  minutely: [
  ],
  hourly: [
    g('coldcaserecords', root),
    g('beta.notify.gov', root),
    g('standards.digital.gov', root),
  ],
  daily: [
    g('digitalcorps.gsa', root),
    g('cloud', root),
    g('oge', root),
    g('search', root),
    wg('cem.va.gov', root),
  ],
  weekly: [
    wg('fac', root),
    wg('usa.gov', root),
  ],
  monthly: [
  ],
  quarterly: [
    g('gsa', root),
  ],
}
