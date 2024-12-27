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

## considerations

It would be nice to be able to uniquely identify:

1. A domain
2. A subdomain
3. A top-level path
4. A path

in a single value. That is, to have a single integer value that is structured, sortable/filterable, and provides information all the way down to the path level. 

* We do not need 255 TLDs; we might need fewer than 16 (2^4, or one nibble)
* We do not need 16M domains. We need fewer than 100K. 2^16 is 65K (four nibbles), 2^18 is 262K.
* We do not need 16M subdomains under every domain. It is probably fewer than 4K (three nibbles).

This suggests 

F:FFFF:FFF:...

meaning we have 8x4=32 bits, and therefore half of the number remains for representing paths.

Across 215K paths, we have 4300 unique path roots. For a given domain space, it might be in the hundreds.

Using those 32 bits for paths, we could:

* Use the first two bits to indicate how we are using the path. 
    * 00 means no structure; treat the next 30 bits (1B) as unique paths
    * 01 means we used two nibbles for the root (64 roots) and 6 nibbles for paths (16M)
    * 10 means we used three nibbles for the root (1024 roots) and 5 nibbles for paths (1M)
    * 11 is undefined 

This would make subpath searching optimal. We can filter, based on the domain64, down to the path

Knowing if we can do this a priori is the trick; that is, what path structure is appropriate for a given site? It might be that we have to assume `00`, and then under analysis (post-crawl), potentially re-assign, which allows for optimization after a second crawl?

Or, we ask our partners.