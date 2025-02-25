# {
#   tld: 1, // a key into tlds.libsonnet
#   id: 1,
#   domain: "18f",

#   // We must explicitly encode the root with the magic key `_root`. 
#   subdomains: {
#     _root: 0,
#     www: 1,
#     acquisition: 2,
#   },

#   // We will use Go-style camel-case.
#   indexFrequency: weekly,

#   // We can imagine adding per-domain 
#   // configuration as needed to this design.
#   ignorePaths: {
#     "*":   ["/something", "*/another_thing"],
#     "acquisition": ["/lots-of-pdfs", "images"],
#   }

#   // There might be paths we *only* crawl.
#   onlyPaths: {
#       "_root": ["/magic"],
#       "acquisition": ["/savings", "/agile"],
#   }
# }

from structures import (
    SiteScannerParser,
    andmap
)
from helpers import (
    subdomain_tree,
)
from bigtree import (
    Node,
    utils,
    find_name
)
import sqlite3

def create_db():
    db_file = "sites.sqlite"
    conn = sqlite3.connect(db_file)
    return conn

# The process will be to build up a new table, one step at a time
# that lets this structure be generated easily.

def create_tld_table(conn):
    q1 = """
    CREATE TABLE IF NOT EXISTS tlds
    (
        id INTEGER PRIMARY KEY,
        tld TEXT UNIQUE
    )
    """
    conn.execute(q1)
    conn.commit()

    q2 = """
    INSERT OR IGNORE INTO tlds (id, tld) VALUES 
    (1, 'gov'), 
    (2, 'mil'),
    (3, 'edu'),
    (4, 'com'),
    (5, 'net'),
    (6, 'org'),
    (7, 'us')
    """
    conn.execute(q2)
    conn.commit()
    return conn

def create_domain_table(conn):
    q = """
    CREATE TABLE IF NOT EXISTS domains
    (
        -- The domain_id is the PK and also our unique D64 value
        domain_id INTEGER PRIMARY KEY,
        tld NUMBER REFERENCES tlds(id),
        -- Never have the same domain appear twice; skip it
        domain TEXT UNIQUE ON CONFLICT ignore,
        index_frequency TEXT,
        source TEXT
    )
    """
    conn.execute(q)
    conn.commit()
    return conn

def is_valid_tld(conn, tld):
    q = """
    SELECT EXISTS(SELECT tld FROM tlds where tld = ?)
    """
    res = conn.execute(q, [tld])
    v = res.fetchone()
    return v[0] == 1

def get_all_tlds(conn):
    q = """
    SELECT tld from tlds
    """
    res = conn.execute(q)
    return list(map(lambda tup: tup[0], res.fetchall()))

def get_tld_id(conn, tld):
    q = """
    SELECT id FROM tlds WHERE tld = ?
    """
    res = conn.execute(q, [tld])
    return res.fetchone()[0]

def insert_domain_row(conn, tld, domain, frequency='monthly', source='unknown'):
    q = """
    INSERT INTO domains
    (tld, domain, index_frequency, source)
    VALUES
    (?, ?, ?, ?)
    """
    # print("INSERT", tld, domain, source)
    conn.execute(q, [get_tld_id(conn, tld), domain, frequency, source])
    conn.commit()

def load_sitescanner_domains(conn):
    q = """
    SELECT target_url FROM sitescanner
    """
    res = conn.execute(q)
    for r in res.fetchall():
        url = list(reversed(r[0].split(".")))
        if is_valid_tld(conn, url[0]):
            insert_domain_row(conn, url[0], url[1], source='sitescanner') 

import re

def load_affiliate_domains(conn):
    q = """
    SELECT "Site domains" FROM affiliates
    """
    res = conn.execute(q)
    for r in res.fetchall():
        domains = r[0].split(",")
        # Now we have a list. Get rid of any empty strings
        # This takes care of trailing commas, which yield empties
        noempty = filter(lambda s: len(s)>0, domains)
        # Now, for each, get rid of anything trailing as a path
        nopaths = map(lambda s: re.sub(r"/.*", "", s), noempty)
        # We should now be able to split on the `.`
        for d in nopaths:
            rev = list(reversed(d.split(".")))
            if andmap(lambda d: re.search("[a-z0-9-]+", d), rev):
                if is_valid_tld(conn, rev[0]):
                    insert_domain_row(conn, rev[0], rev[1], source='affiliate') 


def main():
    # First things, build the TLDs table.
    conn = create_db()
    # Load our allowed TLDs
    create_tld_table(conn)
    # Load the domains
    create_domain_table(conn)
    load_sitescanner_domains(conn)
    load_affiliate_domains(conn)

if __name__ in "__main__":
    main()