import sqlite3
import click
import time

@click.command()
@click.argument('db')
@click.argument('query')
def main(db, query):
  sconn = sqlite3.connect(db)
  scur = sconn.cursor()
  start = time.perf_counter()
  print(query)
  scur.execute("select rowid, rank from ionly where content match ? order by rank desc limit 10", [query])
  end = time.perf_counter()
  elapsed = end - start
  print(list(map(lambda p: p[0], scur.fetchall())))
  print(elapsed)

if __name__ in "__main__":
  main()