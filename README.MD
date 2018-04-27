## Description

This is a quick hack to read DC/OS metrics into a Pandas dataframe.

## How to use for DC/OS Metrics

1) Clone this Repo `git clone https://github.com/fabianbaier/dcos-metrics-snapshot.git` and change into its root directory `cd dcos-metrics-snapshot`
2) Have `DCOS CLI` installed and authenticated to your cluster
3) Create an service account and secret `make sa`
4) Run the crawler `dcos marathon app add crawler.marathon.json`
5) Download the `stdout` file into your repos root directory
6) Run `./convert.sh`
7) See it with Panda `python pandapanda.py`

## How to use for Mesos API

1) Clone this Repo `git clone https://github.com/fabianbaier/dcos-metrics-snapshot.git` and change into its root directory `cd dcos-metrics-snapshot`
2) Have `DCOS CLI` installed and authenticated to your cluster
3) Create an service account and secret `make sa`
4) Run the crawler `dcos marathon app add crawlerMesos.marathon.json`
5) Download the `stdout` file into your repos root directory
6) Run `./convert.sh`
7) See it with Panda `python pandapanda.py`

## Disclaimer

Just use with a small cluster
