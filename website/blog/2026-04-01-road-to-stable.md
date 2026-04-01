---
title: Road to 1.0
description: The Minimum Viable Product is done, we're now working towards 1.0
slug: road-to-stable
authors: linkdd
tags: [release, mvp, stable]
---

:tada:
**[FlowG v0.54.0](https://github.com/link-society/flowg/releases/tag/v0.54.0)**
has been released, finalizing the [MVP](https://github.com/link-society/flowg/milestone/2)
milestone of our roadmap.

<!-- truncate -->

## Minimum Viable Product

For the MVP, we wanted to have at least the following features:

 - ability to ingest and store logs
 - interoperability with many third party services, done via
   [Forwarders](https://link-society.github.io/flowg/docs/technical/forwarders)
 - basis of clustering with the most common cluster formation strategies

After months of work, this is now done. Special thanks to our main contributors,
without who we would not be here today:

 - [@glooopsies](https://github.com/glooopsies)
 - [@Minnerlas](https://github.com/Minnerlas)
 - [@atpugtihsrah](https://github.com/atpugtihsrah)
 - [@n4vxn](https://github.com/n4vxn)

**FlowG** has all you need to make sense of the logs of your production
environments, and we're quite active on the
[bug tracker](https://github.com/link-society/flowg/issues), so don't hesitate
to ask any question, or report eventual problems.

There is still a lot of work to be done, and it's all going towards the **1.0**,
first "stable" version.

## Road to 1.0

What we call "stable" here means that the API will be locked down, no breaking
change guaranteed.

We have 2 big projects for this, and probably a few more down the line that are
not planned yet, or not even fleshed out yet:

### Frontend Redesign

Historically, the frontend was made with [Templ](https://templ.guide),
[HTMX](https://htmx.org), and [Tailwind](https://tailwindcss.com). Due to the
difficulties of handling the amount of client-state we had, and keeping the API
and UI features in sync, we decided to migrate to a [React](https://react.dev)
SPA with [React MUI](https://mui.com).

Unfortunately, lots of Tailwind code remained, and the result is a mix of 2
conflicting design systems, leading to some visual glitches and inconsistencies
that are hard to debug and/or fix.

A proper redesign is waranted:

 - removing Tailwind and fully adopting MUI
 - reorganizing the code to be more consistent and following best practices
 - documenting the architecture to facilitate onboarding new contributors

Over the next few months, this work is gonna be done by
[@coyote2190](https://github.com/coyote2190), our most recent contributor.
Special thanks to him as well!

### Replication

We currently have an experimental implementation of the replication layer,
allowing **FlowG** to be a highly available service.

This implementation is buggy, costly in terms of performance, and wrong on many
levels.

A complete redesign is also planned for the next few months. We'll keep using
the [SWIM Protocol](https://en.wikipedia.org/wiki/SWIM_Protocol), via
[Hashicorp Memberlist](https://github.com/hashicorp/memberlist) Go package, for
the automatic cluster formation. But instead of relying on the costly
"TCP Push/Pull" mechanism to synchronize storages, we'll implement a proper
Operation Log, and synchronize changes via a custom workflow, separate from the
*SWIM* Protocol.

As for the replication model, we still want
[eventual consistency](https://en.wikipedia.org/wiki/Eventual_consistency), that
hasn't changed. Only the implementation will change, and be correct.

Correctness will be proven with an exhaustive test suite, trying to ensure that
most edge cases are covered.

### Testing the chaos

We want to prove that **FlowG** can perform well in production environments,
where many things can go wrong. To do so, we will continuously improve the test
suite and set up test environments, following the principles of
[Chaos Engineering](https://en.wikipedia.org/wiki/Chaos_engineering).

### Traces and Metrics with OpenTelemetry?

At the moment, **FlowG** only supports logs. The OpenTelemetry integration also
only supports those.

But OpenTelemetry also has "metrics" and "traces", which are essential parts of
any observability solutions. It is not yet clear if we want that in **FlowG** or
not, there is an
[ongoing discussion](https://github.com/link-society/flowg/discussions/595)
about it on Github, and would appreciate any feedback on this.

Nothing has been fleshed out yet for this. It might not be part of the 1.0
release, which is fine because it does not look like that introducing this
feature would be a breaking change.

## Conclusion

The list of contributors keeps growing, and you have all our thanks for that!

Any feedback is welcomed, feel free to join our community either on the Github
or on [Discord](https://discord.com/invite/zjG3mMaENg).

If you are using **FlowG** in production, give us a ping, we'll add you to our
`README`.

> To infinity and beyond!
