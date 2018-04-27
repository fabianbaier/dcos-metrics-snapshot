import itertools
import json
import pandas as pd

file = "containers.json"

raw_data = json.loads(open(file).read())
data = pd.io.json.json_normalize(raw_data)

print("<3 Happy Panda <3")