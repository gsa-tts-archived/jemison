local hours(n) = n * 60 * 60;
local minutes(n) = n * 60;

local params(top, env, service, params) =
  {
    [top]: {
      [s[1]]: s[2][env]
      for s in params
      if s[0] == service
    },
  };

local parameters = [
  [
    'allowed_hosts',
    {
      cf: 'all',
      container: 'three',
      localhost: 'three',
    },
  ],
];

{
  hours:: hours,
  minutes:: minutes,
  params:: params,
  parameters: parameters,
}
