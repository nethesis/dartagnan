# Athos installation

Instruction for a clean CentOS 7.


1. Install required packages:

   ```
   yum install epel-release centos-release-scl -y
   yum install git wget  -y
   yum install rh-postgresql96 nginx stunnel redis -y
   yum install certbot python2-certbot-dns-digitalocean -y
   ```

2. Clone the GIT repository:

   ```
   git clone https://github.com/nethesis/dartagnan.git
   cd dartagnan
   ```

3. Configure PostgreSQL sevice:

   ```
   scl enable rh-postgresql96 bash
   postgresql-setup --initdb
   systemctl start rh-postgresql96-postgresql
   systemctl enable rh-postgresql96-postgresql
   ```

   Then, copy `pg_hba.conf` to `/var/opt/rh/rh-postgresql96/lib/pgsql/data/`:

   ```
   cp ./roles/athos/files/pg_hba.conf /var/opt/rh/rh-postgresql96/lib/pgsql/data/
   ```

4. Create database:

   ```
   su - postgres -c 'scl enable rh-postgresql96 -- psql  < ./roles/athos/files/database.sql
   ```

   Change the `SECRET` variable with a generated password.
   You can generate a password using:
   ```
   openssl rand -base64 12
   ```

5. Change database authentication:

   Copy `pg_hba.conf` to `/var/opt/rh/rh-postgresql96/lib/pgsql/data/`:

   ```
   cp ./roles/athos/files/pg_hba.conf /var/opt/rh/rh-postgresql96/lib/pgsql/data/
   systemctl restart rh-postgresql96-postgresql
   ```

6. Prepare the working directory and download the files:

   ```
   mkdir -p /opt/dartagnan
   wget https://github.com/nethesis/dartagnan/releases/download/dev/athos -O /opt/dartagnan/athos
   chmod a+x /opt/dartagnan/athos
   cp ./roles/athos/files/athos.service /etc/system/athos.service
   ```

7. Create configuration file, templates, integrations and licenses:

   ```
   cp ./roles/athos/files/config.json /opt/dartagnan/
   cp -r ./athos/templates /opt/dartagnan/
   cp -r ./athos/integrations /opt/dartagnan/
   cp -r ./deploy/roles/athos/files/license_generator.sh /opt/dartagnan/

   ```

   Edit at least the following options inside `/opt/dartagnan/config.json`:

   - `XXXX`, pick a valid country from the list inside `database.sql`
   - `YOUR_DOMAIN` with your server public FQDN
   - db_password
   - auth0 domain and audience
   - Paypal credentials
   - SMTP host and credentials
   - FRAME_URL home page of the portal. The VueJS app will place a login button on top of it

   Following parameters under `notifications` section are used inside the mail templates:

   - portal_title
   - help_url
   - docs_url

8. Start athos daemon:

   ```
   systemctl start athos
   systemctl enable athos
   ```

9. Configure Let's Encrypt with Digital Ocean DNS:

   ```
   mkdir -p /etc/systemd/system/certbot-renew.service.d/
   cp ./roles/athos/files/athos.conf /etc/systemd/system/certbot-renew.service.d/
   echo "dns_digitalocean_token=XXXXXXX" > /etc/letsencrypt/digitalocean.ini
   chmod 600 /etc/letsencrypt/digitalocean.ini
   systemctl enable --now certbot-renew.timer
   certbot certonly --dns-digitalocean --dns-digitalocean-credentials /etc/letsencrypt/digitalocean.ini --dns-digitalocean-propagation-seconds 60 --register-unsafely-without-email --non-interactive --agree-tos --cert-name athos -d YOUR_DOMAIN
   ```

   Replace `XXXXXXX` with your valid token and `YOUR_DOMAIN` with your server public FQDN.

10. Configure nginx

   Given configuration is built to work with Let's Encrypt Certbot

   ```
   setsebool httpd_can_network_connect 1 -P
   cp ./roles/athos/files/ssl.conf /etc/nginx
   cp ./roles/athos/files/virtualhost.conf /etc/nginx/conf.d
   openssl dhparam -out /etc/ssl/certs/dhparam.pem 4096
   system start nginx
   system enable nginx
   ```

   Change `YOUR_DOMAIN` with your domain inside `virtualhost.conf`.

11. Configure the firewall

    Install required packages:
    ```
    yum install iptables-services
    ```

    Copy the firewall configuration:
    ```
    cp ./roles/athos/files/iptables /etc/sysconfig/iptables
    ```

    Apply the configuration:
    ```
    systemctl start iptables
    systemctl enable iptables
    ```

12. Configure redis:

    ```
    cp ./roles/athos/files/redis.conf /etc/redis
    systemctl start redis
    systemctl enable redis
    ```

13. Configure stunnel:

    ```
    cp ./roles/athos/files/stunnel-redis.conf /etc/stunnel/redis.conf
    cp ./roles/athos/files/redis.service /etc/systemd/system
    systemctl start stunnel@redis
    systemctl enable stunnel@redis
    ```

# Aramis installation

1. Build Aramis following istructions from README:
   ```
   https://github.com/nethesis/dartagnan/blob/master/aramis/README.md
   ```

2. Copy the content of `dist` directory inside `/opt/dartagnan/www`

3. Configure [Auth0](https://auth0.com/) account and set `auth0-config.js` accordingly.
