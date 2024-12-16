import json
import os
import random
import re
import psycopg2
from random import randrange
from pprint import pprint


def only_alpha(char):
    # return (65 <= ord(char) <= 90) or (97 <= ord(char) <= 122)
    return 97 <= ord(char) <= 122


def maine():
    unique = {}
    conn = psycopg2.connect(
        database="postgres", user="postgres", password="", host="localhost", port="7654"
    )

    cursor = conn.cursor()
    cursor.execute("SELECT distinct(domain64) from raw_content")
    d64s = cursor.fetchall()

    for i in range(100):
        random_d64 = d64s[randrange(len(d64s) - 1)]
        cursor.execute(
            f"SELECT content from raw_content where tag='body' and domain64={random_d64[0]} order by random() limit 1"
        )
        result = cursor.fetchone()
        words = result[0].split()
        if len(words) > 10:
            for w in range(randrange(len(words) - 1)):
                onlyascii = "".join([s for s in words[w].lower() if only_alpha(s)])
                if 4 < len(onlyascii) < 9:
                    unique[onlyascii] = True

    cursor.close()
    conn.close()
    with open("wordlist.json", "w") as fp:
        json.dump(list(unique.keys()), fp, indent=2)

    #                 str(words[randrange(len(words) - 1)].encode("ascii", errors="ignore"))

    # for filename in databases:
    #     connection = sqlite3.connect(
    #         os.path.join("../../assets/databases/" + filename + ".sqlite")
    #     )
    #     cursor = connection.cursor()
    #     rows = cursor.execute("SELECT text FROM site_index LIMIT 50")
    #     superset = set()
    #     for row in rows:
    #         words = row[0].split(" ")
    #         for word in words:
    #             if len(word) >= 3:
    #                 superset.add(word)

    #     pprint(superset)
    #     ls = list(superset)
    #     # Keep 50 words to choose from
    #     keepers = list()
    #     iter = 0
    #     while len(keepers) < 50:
    #         iter += 1
    #         word = random.choice(ls)
    #         if re.match("^[a-zA-Z0-9]$", word) or (iter > 100):
    #             keepers.append(word)
    #     unique[filename] = keepers
    # with open("wordlists.json", "w") as fp:
    #     json.dump(unique, fp, indent=2)


if __name__ in "__main__":
    maine()
