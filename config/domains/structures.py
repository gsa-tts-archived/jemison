import csv
import re
from collections import defaultdict
from pprint import pprint

# The affiliates list has some domains in the 
# .cn and .tw space. It is unclear if they are domains
# we should index at this time.
# 
# The sitescanner URLs are only gov and mil.
# The affiliates list includes some others.
VALID_TLDS = [
    "gov",
    "mil",
    "edu",
    "com",
    "org",
    "net",
    "info",
    "blog",
    "dev",
    "us",
    "io"
]

class DomainEntry:
    pass

class SiteScannerRow:
    # instance variables
    # _row : int 
    # _values : dict

    def __init__(self, row_number, line):
        self._row = row_number
        self._values = {} 
        for ndx, field in enumerate([
            "target_url",
            "base_domain",
            "top_level_domain",
            "branch",
            "agency",
            "bureau",
            "source_list_federal_domains",
            "source_list_dap",
            "source_list_pulse",
            "source_list_omb_idea",
            "source_list_eotw",
            "source_list_usagov",
            "source_list_gov_man",
            "source_list_uscourts",
            "source_list_oira",
            "source_list_other",
            "source_list_mil_1",
            "source_list_mil_2",
            "omb_idea_public"
        ]):
            self._values[field] = line[ndx]
    def get(self, key):
        return self._values[key]
    
class SiteScannerParser:   
    def __init__(self, csv_filename):
        self._rows = []
        with open(csv_filename, "r") as f:
            reader = csv.reader(f, delimiter=",")
            next(reader)
            for row_number, line in enumerate(reader):
                self._rows.append(SiteScannerRow(row_number, line))


class AffiliateRow:
    def __init__(self, row_number, line):
        self._row = row_number
        self._values = {}
        for ndx, field in enumerate([
            "id",
            "display_name",
            "site_handle",
            "site_domains",
            "search_engine",
            "active",
            "agency",
        ]):
            self._values[field] = line[ndx]
    def get(self, key):
        return self._values[key]

class AffiliateParser:
    def __init__(self, csv_filename):
        self._rows = []
        with open(csv_filename, "r") as f:
            reader = csv.reader(f, delimiter=",")
            next(reader)
            for row_number, line in enumerate(reader):
                # print(line)
                self._rows.append(AffiliateRow(row_number, line))

# These could be
# a.tld
# a.tld,b.tld
# a.tld.b.tld,
# a.b.tld
# a.b.tld/something
# a.b.tld/something/else
# So, to extract the TLD, some work is needed.
def affiliate_to_tlds(aff):
    affiliates = aff.split(",")
    # Now we have a list. Get rid of any empty strings
    # This takes care of trailing commas, which yield empties
    noempty = filter(lambda s: len(s)>0, affiliates)
    # Now, for each, get rid of anything trailing as a path
    nopaths = map(lambda s: re.sub(r"/.*", "", s), noempty)
    # We should now be able to split on the `.`
    tlds = map(lambda s: list(reversed(s.split(".")))[0], nopaths)
    # Get rid of anything not a-z
    onlyaz = filter(lambda s: re.search("[a-z]+", s), tlds)
    # And, only keep valid TLDs.
    valid = filter(lambda s: s in VALID_TLDS, onlyaz)
    return list(valid)

class TLDs:
    tlds = set()
    def __init__(self):
        pass
    def extract(site_domains):
        return affiliate_to_tlds(site_domains)

def andmap(fun, ls):
    return len(list(filter(lambda v: not v, map(fun, ls)))) == 0

# These could be
# a.tld
# a.tld,b.tld
# a.tld.b.tld,
# a.b.tld
# a.b.tld/something
# a.b.tld/something/else
# So, to extract the domain, some work is needed.
class Domains:
    domains = defaultdict(lambda: defaultdict(list))

    def __init__(self):
        pass

    def extract(dls):
        domains = dls.split(",")
        # Now we have a list. Get rid of any empty strings
        # This takes care of trailing commas, which yield empties
        noempty = filter(lambda s: len(s)>0, domains)
        # Now, for each, get rid of anything trailing as a path
        nopaths = map(lambda s: re.sub(r"/.*", "", s), noempty)
        # We should now be able to split on the `.`
        try:
            for d in nopaths:
                rev = list(reversed(d.split(".")))
                if andmap(lambda d: re.search("[a-z0-9-]+", d), rev):
                    return rev
                else:
                    # print("invalid", rev, d)
                    pass
        except:
            pprint("nopaths", list(nopaths)) 