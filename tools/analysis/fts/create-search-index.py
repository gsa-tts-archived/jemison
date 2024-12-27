import sqlite3
import psycopg2
import json

def create_sqlite():
  sconn = sqlite3.connect("fts2.db")
  scur = sconn.cursor()
  scur.execute("""
               create virtual table if not exists
               ionly
               using fts5(content, content='')
               """)
  return (sconn, scur)

def main():
  sconn, scur = create_sqlite()
  # Connect to an existing database
  pconn = psycopg2.connect(dbname="postgres", user="postgres", host="localhost", port=7654)
  # Open a cursor to perform database operations
  pcur = pconn.cursor()

  # Query the database and obtain data as Python objects
  pcur.execute("""SELECT rowid, content 
               FROM searchable_content 
               where domain64 >= %s and domain64 < %s
               ;""",
              [int('0x0100008700000100', 16),
              int('0x0100008700000200', 16)])
  
  fetching = True
  count = 0
  while fetching:
    res = pcur.fetchone()
    if res == None:
      fetching = False
      continue
    prowid, content = res[0], res[1]
    c = scur.execute("""
                 insert into ionly (content) values (?, ?)
                 """, [content])
    count += 1
    if count % 100 == 0:
      print(count)
      sconn.commit()

  sconn.commit()

    

if __name__ in "__main__":
  main()