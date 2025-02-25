from structures import *

affiliates = [
    "9514,Copy of Bureau of Ocean Energy Management,boemgov_searchgov,\"espis.boem.gov,opendata.boem.gov,www.boem.gov\",SearchGov,True,",
]

def test_affiliates():
    r = AffiliateRow(0, affiliates[0])
    pprint(r._values)
    tlds = TLDs.extract(r.get("site_domains"))
    print(tlds)
    ds = Domains.extract(r.get("site_domains"))
    print(ds)