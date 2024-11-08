import json
import os
import random
import re
import sqlite3

from pprint import pprint


databases = [
    "fac.gov",
    "search.gov",
    "cloud.gov",
    "usa.gov",
    "standards.digital.gov",
]


def maine():
    unique = {}
    for filename in databases:
        connection = sqlite3.connect(
            os.path.join("../../assets/databases/" + filename + ".sqlite")
        )
        cursor = connection.cursor()
        rows = cursor.execute("SELECT text FROM site_index LIMIT 50")
        superset = set()
        for row in rows:
            words = row[0].split(" ")
            for word in words:
                if len(word) >= 3:
                    superset.add(word)

        pprint(superset)
        ls = list(superset)
        # Keep 50 words to choose from
        keepers = list()
        iter = 0
        while len(keepers) < 50:
            iter += 1
            word = random.choice(ls)
            if re.match("^[a-zA-Z0-9]$", word) or (iter > 100):
                keepers.append(word)
        unique[filename] = keepers
    with open("wordlists.json", "w") as fp:
        json.dump(unique, fp, indent=2)


if __name__ in "__main__":
    maine()
