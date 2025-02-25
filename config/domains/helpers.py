from structures import (
    TLDs,
    Domains
)

def get_all_tlds(ssp=None, ap=None):
    tlds = []
    if ssp:
        for r in ssp._rows:
            tlds.append(TLDs.extract(r.get("top_level_domain")))
    if ap:
        for r in ap._rows:
            tlds.append(TLDs.extract(r.get("site_domains")))
    compressed = set()
    for tld in tlds:
        if isinstance(tld, list):
            for t in tld:
                compressed.add(t)
        else:
            compressed.add(tld)
    return compressed

def subdomain_tree(ssp=None, ap=None):
    ds = []
    if ssp:
        for r in ssp._rows:
            ds.append(Domains.extract(r.get("target_url")))
    if ap:
        for r in ap._rows:
            ds.append(Domains.extract(r.get("site_domains")))
    # ds is a list of lists.
    return ds
