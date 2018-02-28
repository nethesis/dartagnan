---
title: "Inventory"
permalink: /docs/inventory/
---

Once a day, the inventory collects all machine information using [Facter](https://docs.puppet.com/facter/).
All data are sent to the Dartagnan instance.

The inventory is composted by `facts` Each `fact` is a piece of information generated from a script.
Facter already provides a set of default data like CPU, disks and network hardware information.
Any package can provide custom facts to collect extra data, like the number of configured users, list of installed packages, etc.

Collected data are accessible for each server and can be used to perform usage analysis on active installations.

More techincal stuff at [nethserver-inventory repository](https://github.com/nethserver/nethserver-inventory).
