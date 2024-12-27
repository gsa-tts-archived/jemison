import psycopg2
import sqlite3
import time
from minio import Minio
from minio.error import S3Error

# create table searchable_content (
#   domain64 bigint not null,
#   path text not null,
#   tag text default 'p' not null,
#   content text not null,
#   fts tsvector generated always as (to_tsvector('english', content)) STORED
# )

# 0100008700000100
# 0100008700000200
def main():
  # sconn = sqlite3.connect("fts.db")
  # scur = sconn.cursor()
  # scur.execute("""
  #              create table if not exists
  #              cmap
  #              (prow int, content text)
  #              """)

  # scur.execute("""
  #              create virtual table if not exists
  #              ionly
  #              using fts5(content, content=cmap, content_rowid=rowid)
  #              """)
  # scur.execute("""
  #              create trigger if not exists ionly_insert after insert on cmap begin
  #              insert into ionly (rowid, content) values (new.rowid, new.content);
  #              update cmap set content = null where rowid = new.rowid;
  #              end;
  #              """)
  
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
  pcur.execute("""select 	* from ts_stat('select fts from searchable_content') 
               where 
               length(word) >= %d and length(word) < %d 
               order by nentry desc, ndoc desc
               """, [4, 10])

  fetching = True
  while fetching:
    row = pcur.fetchone()
    if row == None:
      fetching = False
    else:
      # time.sleep(10)
      prowid, content = row
      if prowid % 1000 == 0:
        print(prowid)
        sconn.commit()

      scur.execute("""
                  insert into cmap (prow, content) values(?, ?);
                  """, [prowid, content])
  
  # Close communication with the database
  pcur.close()
  pconn.close()
  
  scur.close()
  sconn.close()


if __name__ in "__main__":
  main()