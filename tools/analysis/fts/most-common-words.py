import json
from pprint import pprint

js = json.load(open("freq-blog-nasa-gov.json"))

s = sorted(js.items(), key=lambda x:x[1])

pprint(s)