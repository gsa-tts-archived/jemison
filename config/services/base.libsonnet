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
    'schedule',
    {
      cf: 'schedules/schedule.json',
      container: 'schedules/three.json',
      localhost: 'schedules/three.json',
    },
  ],
];

{
  hours:: hours,
  minutes:: minutes,
  params:: params,
  parameters: parameters,
}
