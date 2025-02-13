import csv
import json
import _jsonnet
import click
from pprint import pprint
import os
import subprocess
import re

tlds = {}
domains = {}
subdomains = {}
booleans = {}
sources = {}
babs = {}
# ['worldwar1centennial.gov', 'worldwar1centennial.gov', '.gov', 'Executive', 'The United States World War One Centennial Commission', 'The United States World War One Centennial Commission', 'TRUE', 'FALSE', 'TRUE', 'FALSE', 'FALSE', 'FALSE', 'FALSE', 'FALSE', 'FALSE', 'FALSE', 'FALSE', 'FALSE', '']

source_keys = [
    "federal_domains",
    "dap",
    "pulse",
    "omb_idea",
    "eotw",
    "usagov",
    "gov_man",
    "uscourts",
    "oira",
    "other",
    "mil_1",
    "mil_2",
    "omb_idea_public"
]

def fqdn2rfqd(d):
    return list(reversed(d.split(".")))

def get_tld(s):
    return s[1:]

def get_domain(d):
    return "".join(fqdn2rfqd(d)[1])

def get_subdomain(d):
    pieces = d.split(".")
    if len(pieces) == 2:
        return "_root"
    if len(pieces) > 2:
        return ".".join(fqdn2rfqd(d)[2:])

class UnquotedKeysEncoder(json.JSONEncoder):
    def encode(self, o):
        if isinstance(o, dict):
            return "{" + ", ".join(f"{k}: {self.encode(v)}" for k, v in o.items()) + "}"
        return super().encode(o)

@click.command()
@click.argument('csv_file')
def go(csv_file):
    with open(csv_file, "r") as f:
        reader = csv.reader(f, delimiter=",")
        next(reader)
        for i, line in enumerate(reader):
            # build a set of TLDs
            tld = get_tld(line[2])
            root = line[1]
            url = line[0]
            tlds[tld] = True
            domains[get_domain(root)] = tld
            if get_domain(root) in subdomains:
                subdomains[get_domain(url)].append(get_subdomain(url))
            else:
                subdomains[get_domain(url)] = [get_subdomain(url)]
            sources[get_domain(url)] = line[6:]
            babs[get_domain(url)] = line[3:6]

    pprint(tlds)
    pprint(domains)
    pprint(subdomains)

    # Now, we want to output...
    # one file per tld... `gov.gsa.libsonnet`
    # one object per file
    # {
    #   id: 1,
    #   tld: "gov", // a key into tlds.libsonnet
    #   // We must explicitly encode the root with the magic key `_root`. 
    #   subdomains: [
    #     _root: 0,
    #     www: 1,
    #     acquisition: 2,
    #   ],
    #   // We will use Go-style camel-case.
    #   indexFrequency: weekly,
    # }    for k, v in subdomains:
    id_counter = 1
    for dom, kids in subdomains.items():
        obj = {}
        tld = domains[dom]
        obj["id"] = id_counter
        id_counter += 1
        obj["tld"] = tld
        obj["domain"] = dom
        kids_as_dict = {}
        for ndx, k in enumerate(kids):
            kids_as_dict[k] = ndx
        obj["subdomains"] = kids_as_dict
        obj["indexFrequency"] = "weekly"
        obj["branch"] = babs[dom][0]
        obj["agency"] = babs[dom][1]
        obj["bureau"] = babs[dom][2]
        try:
            os.mkdir(tld)
        except:
            pass
        base = tld + "/" + tld + "." + dom
        jsonnetfname = base + ".libsonnet"
        fp = open(jsonnetfname, "w")
        s = json.dumps(obj)
        # WHY?
        # Because jsonnetfmt will not format a document into an indented
        # multiline file unless there is at least one newline inside the object.
        s = s.replace("{", "{\n", 1)
        # Same for the array
        s = s.replace("[", "[\n", 1)
        fp.write(s)
        fp.close()
        subprocess.run(["jsonnetfmt", "-i", "-n", "2", jsonnetfname])
    
    for tld, _ in tlds.items():
        toplvlobj = []
        print(f"processing {tld}")
        fp = open(tld + "/" + tld + ".libsonnet", "w")
        for dom, kids in subdomains.items():
            if domains[dom] == tld:
                safedom = dom.replace("-", "_")
                var = f"{tld}{safedom}"
                fp.write(f"local {var} = import '{tld}.{dom}.libsonnet';\n")
                toplvlobj.append(var)
        fp.write("\n")
        # json.dump(toplvlobj, fp, indent = 2, cls=UnquotedKeysEncoder)
        fp.write("[\n")
        ndx = 0
        for v in toplvlobj:
            if ndx == len(toplvlobj):
                fp.write(f"\t{v}\n")
            else:
                fp.write(f"\t{v},\n")
            ndx += 1
        fp.write("]\n")
        fp.close()






if __name__ in "__main__":
    go()
