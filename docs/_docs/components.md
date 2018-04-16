---
title: "Software components"
permalink: /docs/components/
---

[Dartagnan server](https://github.com/nethesis/dartagnan) is divided in 3 parts:

- **Athos**: REST API server implemented in [Go](https://golang.org/), it handles subcriptions, payments and serve data to Aramis 
- **Aramis**: web interface which uses Athos APIs to display manage machines, alerts and subscriptions
- **Porthos**: stable updates YUM mirror, accessible only from registered machines

## Athos

The **server** REST APIs are documented [here]({{site.api_url}}).

Athos stores all data inside a [PostgreSQL](https://www.postgresql.org/) database.

The API **client** part is shipped by the **nethserver-subscription** RPM and is
made up of three main parts:

- server registration and YUM repositories configuration
- machine information gathering (inventory)
- machine monitoring and alert generation

The client technical documentation can be found at https://github.com/nethserver/nethserver-subscription.

## Aramis

Partial list of features:

- Authentication via social login or ad-hoc account using [Auth0](https://auth0.com/)
- Machine registration with free 30 days trial
- Subscription and payments management via [PayPal](https://www.paypal.com) API
- Alerts management

## Porthos

Porthos serves YUM repositories via HTTPS. Its main features are:

- access with HTTP basic authentication
- upstream repository synchronization
- snapshot views of each repository, to gradually release updates to client tiers

More information about Porthos is available at https://github.com/nethesis/dartagnan/blob/master/porthos/README.md
