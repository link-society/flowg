---
title: Helm Charts!
description: FlowG v0.64.0 introduces Charts to help(m) you deploy faster!
slug: helm-charts
authors: x86ethan
tags: [release, helm, k8s]

---

New FlowG release [0.64.0](https://github.com/link-society/flowg/releases/tag/0.64.0) introduces 3 Helm Charts to enable you to deploy FlowG faster and in an easier way in your own cluster.

## Why Helm?

Since [0.61.0](https://github.com/link-society/flowg/releases/tag/0.61.0) FlowG can run in a cluster using FoundationDB, but up until now, clusterizing it was a bit complicted. With or without Docker, the process relied on many complex and difficult to automate steps. And our guides weren't mentionning Kubernetes, leaving users entirely to their own devices when it cames to setting up a K8s manifest. 

That's why we started working on Helm Charts. Helm provides an easy, portable,  user-friendly, and yet automatable way of setting up any service on a cluster. The process now consists of choosing the architecture that fits your needs, editing the `values.yaml` file, install the chart, and there you go!

## For what needs? 

As I said, there are 3 different charts, each providing a different architecture base on your needs. 

- One provides a **single-node** approach, with only 1 FlowG process, attached to a persistend volume *(PVC)*, which relies on [BadgerDB](https://github.com/dgraph-io/badger) as a back-end.
- Another allows for a **multi-node** setup, this time using [FoundationDB](https://www.foundationdb.org/). 
- The last one follows the same principle as the previous one, while bringing up an additional [Fluentd](http://fluentd.org/) **observability stack** to help you route your cluster's services logs directly into FlowG. 

To learn how to quickly run FlowG on your Kubernetes cluster, just follow [this guide](/docs/user/advanced/k8s/setup.md).

---

## Our architecture decisions 

To bring these charts up, a certain amount of architecture decisions have been made in order to keep a consistent, secure and easy to configure Chart. Let's dive into some of them. 

### Configuring flowg: no variable exposure!!

The documentation about [Flowg's configuration](/docs/technical/configuration) shows a certain number of options you can configure to make FlowG suit your needs. These can be all adjusted via environment variables, first in order to have **easy-to-configure Docker containers**. 

In a pure `docker run` context, or even inside a Docker compose file, it isn't an issue to expose all these environment variables. The user is the only master of the whole keys.  
However, when this environment becomes an auto-generated Kubernetes manifest, errors can easily be made and conflicts can rapidly emerge. 

> Let me give a quick example: what if `FLOWG_HTTP_TLS_CERT` is set to `/etc/ssl/tls.crt`, when the chart injects TLS certs and keys in `/data/ssl/`? The admin will see HTTPS is not working, waste 2 hours figuring out why, and will ultimately have to dig right in the chart's templates, read the 115th line of one of the half-dozen of files that make up the chart, to find a volume mountpoint option, buried in barely-readable Go template options. You understand what I mean?

That's why we made a clear architectural choice that were to drive the development of these charts: [make invalid states unrepresentable](https://www.awwsmm.com/blog/make-invalid-states-unrepresentable). When developing a tool intended for users, leaving nothing to chance is as much a matter of security as it is of ergonomics.

Finally, FlowG's configuration options are fully named in the `values.yaml` file, then written one-by-one in a separate ConfigMap, then given as environment variables to FlowG. This could potentially become a burden if we were to add configuration options to FlowG, but we find it necessary to maintain consistence in the `values.yaml` file. 

