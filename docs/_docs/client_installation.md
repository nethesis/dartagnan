---
title: "Client installation"
permalinks: /docs/client_installation/
---

The [subscription](https://github.com/NethServer/nethserver-subscription) is composed by 3 different parts:

- Subscription: manages the machine registration and configures yum repositories
- Inventory: collects the machine information
- Monitoring: monitors the status of the machine and generates alerts

Currently, installation is supported only on [NethServer 7](https://www.nethserver.org).


## Install on NethServer

1. Install all required packages.

   Access a shell and execute:
   ```
   yum --enablerepo=nethserver-testing install @nethserver-subscription
   ```

2. Access the ``Subscription`` page inside the Server Manager (``http://<SERVER_FQDN>:980``)

3. Follow the istructions on the screen

