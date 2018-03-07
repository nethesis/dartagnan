---
title: "Client installation"
permalinks: /docs/client_installation/
---

The client is composed by 3 parts:

- Subscription: [nethserver-subscription](https://github.com/NethServer/nethserver-subscription), manages the machine registration and configures yum repositories
- Inventory: [nethserver-inventory](https://github.com/NethServer/nethserver-inventory), collects the machine information
- Monitoring: [nethserver-alerts](https://github.com/NethServer/nethserver-alerts), monitors the status of the machine and generates alerts

Currently, installation is supported only on [NethServer 7](https://www.nethserver.org).


## Install on NethServer

1. Install all required packages.

   Access a shell and execute:
   ```
   yum --enablerepo=nethserver-testing install nethserver-subscription-ui nethserver-alerts nethserver-inventory nethserver-yum-cron
   ```

2. Access the ``Subscription`` page inside the Server Manager (``http://<SERVER_FQDN>:980``)

3. Follow the istructions on the screen

