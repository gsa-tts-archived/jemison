local db(module_name, org, space, name, plan, tags=['rds'],) = {
  [module_name]: {
    source: 'github.com/gsa-tts/terraform-cloudgov//database?ref=v0.9.1',
    cf_org_name: org,
    cf_space_name: space,
    name: name,
    recursive_delete: false,
    tags: tags,
    rds_plan_name: plan,
  },
};

local _s3(org, space, name, plan, tags=['s3']) = {
  source: 'github.com/gsa-tts/terraform-cloudgov//s3?ref=v0.9.1',
  cf_org_name: org,
  cf_space_name: space,
  name: name,
  recursive_delete: false,
  tags: tags,
  s3_plan_name: plan,
};

local s3(module_name, org, space, name, plan, tags=['s3']) = {
  [module_name]: _s3(org, space, name, plan, tags=tags),
};

local s3_multi(names, org, space, plan, tags=['s3']) = {
  ['s3-private-' + name]: _s3(org, space, name, plan, tags=tags)
  for name in names
};

local _env(ndri) = if std.objectHas(ndri, 'env')
then { environment: ndri.env }
else {};

local _routes(ndri) = if std.objectHas(ndri, 'routes') && std.length(ndri.routes) > 0
          then { routes: [{ route: id} for id in ndri.routes] }
          else {};

local _bindings(ndri) = if std.objectHas(ndri, 'bindings') && std.length(ndri.bindings) > 0
     then { service_binding: [
       { service_instance: b }
       for b in ndri.bindings
     ] }
     else {};

local _app(ndri, timeout=200, health_check_timeout=180) = {
    name: ndri.name,
    space: 'data.cloudfoundry_space.app_space.id',
    buildpacks: ['https://github.com/cloudfoundry/apt-buildpack', 'https://github.com/cloudfoundry/binary-buildpack.git'],
    path: 'zips/' + ndri.name + '.zip',
    source_code_hash: 'filesha256("zips/' + ndri.name + '.zip")',
    disk_quota: ndri.disk,
    memory: ndri.ram,
    instances: ndri.instances,
    strategy: 'rolling',
    timeout: 200,
    health_check_type: 'port',
    health_check_timeout: health_check_timeout,
    health_check_http_endpoint: '/heartbeat',
} + _env(ndri) + _routes(ndri) + _bindings(ndri);

/*
"service_binding": [
        {[b.name]: { service_instance: "module." + b.id} for b in bindings},
      ],
      "routes": [{route: "cloudfoundry_route." + id} for id in routes],
      },
      */

// local service(name, disk, ram, instances, timeout=200, health_check_timeout=180) = {
//   resource: {
//     cloudfoundry_app:
//       _app({ name: name, disk: disk, ram: ram, instances: instances }, env={}, timeout=200, health_check_timeout=180),
//     service_binding: [
//       { [b.name]: { service_instance: 'module.' + b.id } for b in bindings },
//     ],
//     routes: [{ route: 'cloudfoundry_route.' + id } for id in routes],
//   },
// };

local services(ndris) = {
    cloudfoundry_app: {
      [ndri.name]: _app(ndri)
      for ndri in ndris
    },
};

{
  db:: db,
  s3:: s3,
  s3_multi:: s3_multi,
  services:: services,
}
