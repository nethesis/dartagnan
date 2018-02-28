---
title: "Software components"
permalink: /docs/components/
---

[Dartagnan server](https://github.com/nethesis/dartagnan) is divided in 4 parts:

- **Athos**: REST API server implemented in [Go](https://golang.org/), it handles subcriptions, payments and serve data to Aramis 
- **Aramis**: web interface which uses Athos APIs to display manage machines, alerts and subscriptions
- **Porthos**: stable updates YUM mirror, accessible only from registered machines

Athos stores all data inside a [PostgreSQL](https://www.postgresql.org/) database.

The client part is composed by:

- [nethserver-subscription](https://github.com/nethserver/nethserver-subscription): server registration and YUM repositories configuration
- [nethserver-inventory](https://github.com/nethserver/nethserver-inventory): machine information gathering
- [nethserver-alerts](https://github.com/nethserver/nethserver-alerts): machine monitoring and alert generation

All techincal documentation can be found in the links above.

## Athos

All rest APIs are available [here](https://documenter.getpostman.com/view/3364668/dartagnan/RVfyBA1b#6014623b-7bf6-82ba-d328-8f4d18e76cb5).

## Aramis

Partial list of featues:

- Authentication via social login or ad-hoc account using [Auth0](https://auth0.com/)
- Machine registration with free 30 days trial
- Subscription and payments management via [PayPal](https://www.paypal.com) API
- Alerts management

