---
title: Using local cloud emulators
description: FlowG v0.45.0 introduces API compatibility with ElasticSearch
slug: using-floci-local-emulators
authors: gloopsies
tags: []
---

Say welcome to Amazon AWS, Google Cloud and Microsoft Azure!

<!-- truncate -->

With the release of FlowG 1.0 [steadily approaching](/blog/road-to-stable), our team looked into what our users could be missing from our first stable release.
While already having eight general-purpose forwarders, we were missing direct integrations with major cloud providers. 

All three of the biggest cloud providers have their own solutions for collecting and managing logs, with proprietary implementations and custom libraries.

While using the libraries is generally simple, testing them proved to be much harder.
Testing software integrations with major cloud providers usually requires accounts for each service, along with significant costs if they’re part of a CI/CD pipeline that runs regularly to ensure no breaking changes are introduced to the project.

## Using floci for AWS CloudWatch emulation

The first integration we implemented was for the AWS CloudWatch service.
While searching for local cloud emulators, [floci](https://floci.io/) stood out as the most professional-looking solution.

Using floci couldn't be simpler. All that was required was a single Docker command, and you can have a local version of many AWS services.

```shell
docker run -p 4566:4566 floci/floci:latest
```

Configuring the go `aws-sdk-go-v2` library to use a local endpoint is even easier: just pass a simple url string to the configuration struct, and you're using a local floci instance.

We could create log groups and streams with the official AWS CLI tool or Python library and see forwarded logs working as expected in no time.

```shell
export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_DEFAULT_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test

aws logs create-log-group  --log-group-name flowg
aws logs create-log-stream --log-group-name flowg --log-stream-name logs
```

## Looking for Google Cloud Logging emulator

Full of hope after the AWS success story, we expected that adding Google Cloud Logging would be just as simple.
There were a few solutions we found after a bit of research that looked very similar to what floci provided for AWS, but after trying all of them, none of the solutions actually worked. 

Some didn't even have functional download links, or the instructions provided were completely wrong.

While we researched other solutions, floci team quietly worked on [floci-gcp](https://floci.io/gcp/) and added features we were missing earlier.

Realizing we could use floci again, we were greeted by its simplicity and ease of use once more. While Google's library required a bit of configuration to allow insecure local connections, it didn't require any stream or group creation and was just as simple to use as AWS.  

## Final boss: Microsoft Azure

Since other cloud providers were this simple, and we realized that [floci-az](https://floci.io/az/) already supported Azure Monitor, we thought that this would be just the routine and would be done in no time.

But Microsoft always finds a way to complicate things.

All our previous local cloud servers used HTTP since there is no authentication and they're only used for CI/CD, not actual deployments.
AWS works without changes, Google requires explicit configuration for HTTP, and so does Microsoft, until you actually try to use it.

Microsoft does allow us to set the `InsecureAllowCredentialWithHTTP` flag, but after some debugging, it became apparent that the flag was never read.

Thankfully `floci-az` helped us there as well, with option to run the server over https with the `FLOCI_AZ_TLS_ENABLED` environment variable.

Microsoft also required us to create a custom authentication token implementation if we didn't want to use the system-wide configuration, and we had to use two different client libraries to create the workspace, workgroup, and DCR needed to send the logs.

Unless your infrastructure already heavily depends on Microsoft Azure, I wouldn't recommend using it for your projects.

## Authentication

The biggest advantage of using floci is also its greatest weakness.

Floci doesn't require any authentication and works with any random string instead of a secure token.
This is great for testing functionality and using it inside a CI/CD pipeline, but it doesn't give us the confidence that our integration works with the actual provider's infrastructure.

That's why, at the end of the day, we still need to try every cloud provider manually at least once and check that our authentication actually works in the real world.

While floci is great and has helped us a lot on this journey, sometimes you have to do things yourself.