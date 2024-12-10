// WARNING
// Renumbering these once everything is running/deployed
// will break many tables. We count on these for efficient
// storage of this data. If a content type (e.g.) changes
// value, the application logic breaks everywhere.
// Once encoded, do not renumber.
local schemes = {
  https: 1,
  http: 2,
  sftp: 3,
  ftp: 4,
};

local tlds = {
  gov: 1,
  mil: 2,
  com: 3,
  edu: 4,
  net: 5,
  org: 6,
};

local content_types = {
  'binary/octet_stream': 1,
  'text/html': 2,
  'application/pdf': 3,
};

local tags = {
  unknown: 1,
  title: 2,
  p: 3,
  div: 4,
  h1: 5,
  h2: 6,
  h3: 7,
  h4: 8,
  h5: 9,
  h6: 10,
  a: 11,
  th: 12,
  td: 13,
  page: 14,
};

// These cannot conflict with the values in tags.
local attributes = {
  'img.alt': 1000,
};

local flipMap(m) = { [std.toString(m[k])]: k for k in std.objectFields(m) };

{
  SchemeToConst: schemes,
  ConstToScheme: flipMap(schemes),
  TldToConst: tlds,
  ConstToTld: flipMap(tlds),
  ContentTypeToConst: content_types,
  ConstToContentType: flipMap(content_types),
  TagToConst: tags,
  ConstToTag: flipMap(tags),
  AttributeToConst: attributes,
  ConstToAttribute: flipMap(attributes),
  // We don't want the tags and attributes to overlap
  // as "namespaces." This is because we need to be able to 
  // differentiate them for ranking purposes in SQL.
  assert std.setInter(
    std.set(std.objectValues(self.TagToConst)),
    std.set(std.objectValues(self.AttributeToConst))
  ) == [],
}
