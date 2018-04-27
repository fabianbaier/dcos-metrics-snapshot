import itertools
import json
import pandas as pd

file = "containers.json"
raw_data = json.loads(open(file).read())

# For DC/OS Metrics
# data = pd.io.json.json_normalize(raw_data)

# For Mesos
# for i in range(len(raw_data)):
#    data = pd.io.json.json_normalize(raw_data[i])

print("<3 Happy Panda <3")