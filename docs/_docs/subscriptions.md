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

Stable updates repositories are clones of NethServer and upstream YUM
repositories. Any RPM update is delayed for a week at least after it has been
released by upstream, then it is gradually released to all the subscribed machines
using `tiers`.

A `tier` is a snapshot of a repository in a given time and updates flow from the
tier 0 (`t0`) to tier 3 (`t3`).

General rules:

* Each machine is automatically associated to a tier.
* Tiers are considered only by overnight/automated updates
* Day-time/manual updates from system console or Software Center page always 
  point to the base tier

#### When a specific RPM will be available on stable updates repositories?

If the RPM follows the automated schedule, and is released during *week 0* by
upstream repositories, then it is available during *week 2* from stable
repositories.

In other words, the **minimum age** of an RPM from stable repositories is **one
week**.

This is the automated overnight updates schedule:

| tier | day |
|---|---|
| t0 | Tue |
| t1 | Wed |
| t2 | Thu |
| t3 | Fri |
