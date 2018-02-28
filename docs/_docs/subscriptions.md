---
title: "Subscriptions"
permalink: /docs/subscriptions/
---

A subscription is an entitlement associated to a machine.
Each machine is uniquely identified with a `System ID` and a `Secret`.

The `System ID` is the name fol the machine registered at Dartagnan instance, while
the `Secret` is mainly used to invoke remote APIs to retrieve subscription plans information
and send inventory along with alerts.

Both credentials are used also to access stable updates repositories using HTTP basic authentication.

## Machine registration

`System ID` and `Secret` are generated from the Dartagnan instance and must be saved inside
client, see [nethserver-subscription](https://github.com/NethServer/nethserver-subscription) for more technical details.

## Stable updates repositories

Stable updates repositories are clones of NethServer and upstream YUM repositories.
All updates are delayed for a week after they have been released from upstream,
then gradually released to all subscribed machines using `tiers`.

A `tier` is a snapshot of a repository in a given time and updates flow from the `tier0` to `tier4`.

General rules:

* Each machine is automatically associated to a `tier`.
* Tiers are considered only by overnight/automated updates
* Day-time/manual updates from system console or Software Center page always 
  point to the base tier

#### When a specific RPM will be available on stable updates repositories?

If the RPM follows the automated schedule, and is released during *week 1* in upstream repositories,
then it is available during *week 2* in stable repositories. This is the automated overnight updates schedule:

| tier | day |
|---|---|
| t0 | Tue |
| t1 | Wed |
| t2 | Thu |
| t3 | Fri |
