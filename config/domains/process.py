from structures import (
    SiteScannerParser,
    AffiliateParser,
    TLDs
)

def get_all_tlds(ssp, ap):
    for r in ssp._rows:
        TLDs.add(r.get("top_level_domain"))
    for r in ap._rows:
        TLDs.add(r.get("site_domains"))

def main():
    ssp = SiteScannerParser("sitescanner.csv")
    ap = AffiliateParser("affiliates.csv")
    # Statefully builds up the value
    # in the TLDs class.
    get_all_tlds(ssp, ap)

if __name__ in "__main__":
    main()