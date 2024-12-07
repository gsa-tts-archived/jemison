# domain64

`domain64` is a BIGINT (or 64-bit) type that can be used to encode all domains we are likely to encounter. It represents well as JSonnet/JSON, and can be used in partitioning database tables easily.

## what is it


To encode all of the TLDs, domains, and subdomains we will encounter, we'll use a `domain64` encoding. It maps the entire URL space into a single, 64-bit number (or, `BIGINT`).

```
FF:FFFFFF:FFFFFF:FF
```

or

```
tld:domain:subdomain:RESERVED
```

This lets us track 

* 255 (#FF) TLDs 
* 16,777,215 (#FFFFFF) domains under each TLD
* 16,777,215 (#FFFFFF) subdomains under each domain
* 255 (#FF) reserved

For example

```
01:000001:000000:00 gov.gsa
01:000001:000001:00 gov.gsa.tts
01:000001:000002:00 gov.gsa.api.staging
01:000001:000003:00 gov.gsa.api.prod

```

Subdomains might be deeper: gov.nasa.`earthdata.uat.cdn`, for example. 

| tld | domain | sub | hex | dec |
| --- | --- | --- | --- | --- |
| gov | gsa |  _  | #x0100000100000000   | 72057598332895232 |
| gov | gsa | tts | #x0100000100000100   | 72057598332895488 |
| gov | gsa | api | #x0100000100000200   | 72057598332895744 |

## for partitioning

On a table that contains a `domain64` value, we can partition based on numeric ranges very efficiently.


```sql
CREATE TABLE govgsa PARTITION OF gov
    FOR VALUES FROM (72057594054705152) TO (72057594071482367);
```

Or

```sql
CREATE TABLE govgsatts PARTITION OF gov
    FOR VALUES FROM (72057594054705408) TO (72057594054705663);
```

## As Jsonnet/JSON

Jsonnet will naturally sort by the hex key values.

```
{
  "01": {
    "name": "gov",
    "children": {
      "00000001": {
        "name": "gsa",
        "children": {
          "0001": "tts",
          "0002": "18f",
          "0003": "api"
        }
      }
    }
  },
  "02": {
    "name": "mil",
    "children": {
      "00000001": {
        "name": "af",
        "children": {
          "0001": "www"
        }
      }
      "00000002": {
        "name": "navy"
      }
    }
  }
}
```