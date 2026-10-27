---
title: Helm Charts!
description: FlowG v0.64.0 introduces Charts to help(m) you deploy faster!
slug: helm-charts
authors: x86ethan
tags: [release, helm, k8s]
---

# Helm Charts!

New FlowG release [0.64.0](https://github.com/link-society/flowg/releases/tag/0.64.0) introduces 3 Helm Charts to enable you to deploy FlowG faster and in an easier way in your own cluster.

# Why Helm?

Since [0.61.0](https://github.com/link-society/flowg/releases/tag/0.61.0) FlowG can run in a cluster using FoundationDB, but up until now, clusterizing it was a bit complicted. With or without Docker, the process relied on many complex and difficult to automate steps. And our guides weren't mentionning Kubernetes, leaving users entirely to their own devices when it cames to setting up a K8s manifest. 

That's why we started working on Helm Charts. Helm provides an easy, portable,  user-friendly, and yet automatable way of setting up any service on a cluster. The process now consists of choosing the architecture that fits your needs, editing the `values.yaml` file, install the chart, and there you go!

# For what needs? 

As I said, there are 3 different charts, each providing a different architecture base on your needs. 

- One provides a **single-node** approach, with only 1 FlowG process, attached to a persistend volume *(PVC)*, which relies on [BadgerDB](https://github.com/dgraph-io/badger) as a back-end.
- Another allows for a **multi-node** setup, this time using [FoundationDB](https://www.foundationdb.org/). 
- The last one follows the same principle as the previous one, while bringing up an additional [Fluentd](http://fluentd.org/) **observability stack** to help you route your cluster's services logs directly into FlowG. 