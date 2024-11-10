import sqlite3
import click
import os
import requests
import random
import pathlib
from datetime import timedelta, datetime
import hashlib

word_site = "https://www.mit.edu/~ecprice/wordlist.10000"

response = requests.get(word_site)
WORDS = [w.decode("utf-8") for w in response.content.splitlines()]

t1 = """
CREATE TABLE IF NOT EXISTS
schemes (
  scheme TEXT
);
"""

t2 = """
CREATE TABLE IF NOT EXISTS
hosts (
  host TEXT
);
"""

t3 = """
CREATE TABLE IF NOT EXISTS
metadata (
  scheme INTEGER,
  host INTEGER,
  path TEXT,
  last_updated DATE,
  last_fetched DATE,
  sha1 TEXT,
  length INTEGER,
  content_type TEXT,
  -- constraints must come last
  FOREIGN KEY(scheme) REFERENCES schemes(ROWID),
  FOREIGN KEY(host) REFERENCES hosts(ROWID)
);
"""

tables = [t1, t2, t3]

hosts = [f"{host}.gov" for host in [random.choice(WORDS) for _ in range(200)]]


def SHA1(msg):
    return hashlib.sha1(msg.encode()).hexdigest()


def gen_datetime(min_year=1900, max_year=datetime.now().year):
    # generate a datetime in format yyyy-mm-dd hh:mm:ss.000000
    start = datetime(min_year, 1, 1, 00, 00, 00)
    years = max_year - min_year + 1
    end = start + timedelta(days=365 * years)
    return start + (end - start) * random.random()


@click.command()
@click.argument("rows")
def main(rows):
    db_filename = pathlib.Path(f"test-{rows}.sqlite")
    db_filename.unlink(missing_ok=True)
    connection = sqlite3.connect(db_filename)
    cursor = connection.cursor()
    for t in tables:
        cursor.execute(t)
    for s in ["http", "https"]:
        cursor.execute(f"INSERT INTO schemes VALUES ('{s}')")
    for h in hosts:
        cursor.execute(f"INSERT INTO hosts VALUES ('{h}')")
    connection.commit()
    for n in range(int(rows)):
        if (n > 0) and (n % 10000) == 0:
            connection.commit()
            print(f"Inserted {n} rows...")
        host = random.randrange(1, len(hosts) + 1)
        path = "-".join([random.choice(WORDS) for n in range(random.randrange(3, 8))])
        sha = SHA1(f"{host}/{path}")
        cursor.execute(
            "INSERT INTO metadata VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
            [
                random.choice([1, 2]),
                host,
                path,
                gen_datetime(),
                gen_datetime(),
                sha,
                random.randrange(3_000, 25_000_000),
                random.choice(WORDS) + "/" + random.choice(WORDS),
            ],
        )
    connection.execute("CREATE INDEX metadata_host_path_ndx ON metadata(host, path)")
    connection.execute("CREATE INDEX metadata_lf_ndx ON metadata(last_fetched)")
    connection.execute("CREATE INDEX metadata_ct_ndx ON metadata(content_type)")
    connection.commit()
    connection.execute("ANALYZE")
    connection.execute("VACUUM")
    connection.commit()
    connection.close()


if __name__ in "__main__":
    main()
