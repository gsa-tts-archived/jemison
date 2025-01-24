# file_uploader.py MinIO Python SDK example
from minio import Minio
from minio.error import S3Error
import time
from bs4 import BeautifulSoup
import re
import string
import sqlite3
import sys
import psycopg2
import json
from pprint import pprint
from collections import defaultdict
from unidecode import unidecode
import demoji

def remove_html_tags(text):
    soup = BeautifulSoup(text, 'html.parser')
    clean_text = soup.get_text()
    return clean_text

def process(txt, most_common_words):
  clean = []
  txt = txt.lower()
  txt = demoji.replace(txt, " ")

  for c in ["’", "”", "“", "‘", "!", "…", "-", "—", "®"]:
    txt = txt.replace(c, " ")
  
  for c in ["\n", "\r", "\t"]:
    txt = txt.replace(c, " ")

  txt = txt.replace("\xa0", " ")
  txt = txt.replace("'s", "")
  
  for c in string.punctuation:
    txt = txt.replace(c, " ")

  for s in txt.split(' '):  
    if re.search('[0-9]', s):
       continue
    if len(s) <= 3:
       continue
    if len(s) > 12:
       continue
    if re.search('http', s):
       continue
  
    # https://stackoverflow.com/questions/517923/what-is-the-best-way-to-remove-accents-normalize-in-a-python-unicode-string
    s = unidecode(s)
    s = s.strip()
    # if s in most_common_words and most_common_words[s] > 9100:
    #   pass
    # else:
    clean.append(s)
  # print(clean)
  return clean
  
def get_raw_object_content(mc, obj):
  content = "could not get content"
  name = obj.object_name
  if name.endswith("json"):
    name = obj.object_name.replace("json", "raw")
  try:
    response = mc.get_object('fetch', name)
    raw = response.read().decode('utf-8')
    content = remove_html_tags(raw)
  except Exception as e:
    print("fail:", obj.object_name)
    print(e)
    response.close()
    response.release_conn()
    return None
  finally:
      response.close()
      response.release_conn()
  return content


def get_json_object(mc, obj):
  name = obj.object_name
  if name.endswith("raw"):
    name = obj.object_name.replace("raw", "json")
  try:
    response = mc.get_object('fetch', obj.object_name)
    raw = response.read().decode('utf-8')
    json_obj = json.loads(raw)
  except Exception as e:
    print("fail:", obj.object_name)
    print(e)
    # print(response.status)
    # time.sleep(0.5)
    response.close()
    response.release_conn()
    return None
  finally:
      response.close()
      response.release_conn()
  return json_obj

def is_acceptable_path(p):
  bad = [r"/category/", r"right.*here", r"/tag/",
         r"/author/", r"/page/", r"wp-login", r"wp-content",
         r"/\d{4}/\d{2}$"]
  for b in bad:
    if re.search(b, p):
      return False
  return True

def main():
  mc = Minio("localhost:9100",
          access_key="numbernine",
          secret_key="numbernine",
          secure=False,
      )
  bucket_name = "fetch"
  pconn = psycopg2.connect(dbname="postgres", 
                           user="postgres", 
                           host="localhost", 
                           port=7654)
  # Open a cursor to perform database operations
  pcur = pconn.cursor()

  # Make the bucket if it doesn't exist.
  found = mc.bucket_exists(bucket_name)
  if found:
    print("yay")

  objects = mc.list_objects('fetch', 
                            prefix="blogs.nasa.gov",
                            recursive=True)
  cnt = 0
  freq = defaultdict(int)

  most_common_words = json.load(open("freq-blog-nasa-gov.json"))


  for obj in filter(lambda o: "json" in o.object_name, objects):
    content = None
    metadata = None

    content = get_raw_object_content(mc, obj)
    metadata = get_json_object(mc, obj)

    # print("content", content)
    # print("meta", metadata)
    # print("-----------------------------")
    
    if content != None and metadata != None:
      if is_acceptable_path(metadata['path']):
        cnt += 1
        if cnt % 100 == 0:
          # Get data of an object.
          print(cnt)
          pconn.commit()        
        clean = process(content, most_common_words)
        for word in clean:
          freq[word] += 1

        # print(clean)
        r = pcur.execute("insert into raw2 (domain64, path, tag, content) values (%s, %s, %s, %s)", 
                      [
                        int("0x0100008700000100", 16),
                        metadata["path"],
                        'body',
                        " ".join(clean)
                        ])
  print(cnt)
  json.dump(freq, open("freq-blog-nasa-gov.json", "w"))

if __name__ in "__main__":
  main()