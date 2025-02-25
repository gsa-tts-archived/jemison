from structures import (
    SiteScannerParser,
    AffiliateParser,
)

from helpers import (
    get_all_tlds,
    subdomain_tree,
)

from bigtree import (
    Node,
    find_name,
    utils
)

def main():
    ssp = SiteScannerParser("sitescanner.csv")
    ap = AffiliateParser("affiliates.csv")
    # Statefully builds up the value
    # in the TLDs class.
    all_tlds = get_all_tlds(ssp, ap)
    roots = {}
    for tld in all_tlds:
        roots[tld] = Node(tld)
    subs = subdomain_tree(ssp, ap)
    for sub in list(filter(lambda s: s and len(s) > 1, subs)):
        # sub is a list version of a RFQDN
        # [gov, blah, blah]
        for ndx, part in enumerate(sub[1:]):
            if sub[0] in roots:
                if ndx == 0:
                    try:
                        Node(part, parent=roots[sub[0]])
                    except:
                        pass
                else:
                    try:
                        # print("parent of", part, find_name(roots[sub[0]], sub[ndx]))
                        Node(part, parent=find_name(roots[sub[0]], sub[ndx]))
                    except:
                        pass
    for tld, root in roots.items():
        root.show()
        count = 0
        last = None
        for last in utils.iterators.preorder_iter(root):
            if last.children == ():
                count += 1
        print(tld, count)

if __name__ in "__main__":
    main()