# Athos installation

Instruction for a clean CentOS 7.


1. Install required packages:

   ```
   yum install epel-release centos-release-scl -y
   yum install git wget  -y
   yum install rh-postgresql96 nginx -y
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

5. Prepare the working directory and download the files:

   ```
   mkdir -p /opt/dartagnan
   wget https://github.com/nethesis/dartagnan/releases/download/dev/athos -O /opt/dartagnan/athos
   chmod a+x /opt/dartagnan/athos
   cp ./roles/athos/files/athos.service /etc/system/athos.service
   ```

6. Create a configuration file:

   ```
   cp ./roles/athos/files/config.json /opt/dartagnan/
   
   ```

   Edit at least the following options inside `/opt/dartagnan/config.json`:

   - `YOUR_DOMAIN` with your server public FQDN
   - db_password
   - auth0 domain
   - Paypal credentials

7. Start athos daemon:

   ```
   systemctl start athos
   systemctl enable athos
   ```

8. Configure Let's Encrypt with Digital Ocean DNS:

   ```
   mkdir -p /etc/letsencrypt/
   echo -e '#!/bin/bash\nsystemctl restart nginx' > /usr/local/bin/letsencrypt-hook.sh
   chmod a+x /usr/local/bin/letsencrypt-hook.sh
   echo "dns_digitalocean_token=XXXXXXX" > /etc/letsencrypt/digitalocean.ini
   chmod 600 /etc/letsencrypt/digitalocean.ini
   echo "certbot certonly --dns-digitalocean --dns-digitalocean-credentials /etc/letsencrypt/digitalocean.ini --dns-digitalocean-propagation-seconds 60 --register-unsafely-without-email --non-interactive --agree-tos --post-hook letsencrypt-hook.sh  -d YOUR_DOMAIN" > /etc/cron.daily/letsencrypt
   chmod a+x /etc/cron.daily/letsencrypt
   ```
 
   Replace `XXXXXXX` with your valid token and `YOUR_DOMAIN` with your server public FQDN.

9. Configure nginx

   Given configuration is built to work with Let's Encrypt Certbot

   ```
   setsebool httpd_can_network_connect 1 -P
   cp ./roles/athos/files/ssl.conf /etc/nginx
   cp ./roles/athos/files/virtualhost.conf /etc/nginx/conf.d
   openssl dhparam -out /etc/ssl/certs/dhparam.pem 4096
   system start nginx
   system enable nginx
   ```

   Change YOUR_DOMAIN` with your domain inside `ssl.conf` and `virtualhost.conf`.

# Aramis installation

1. Build Aramis following istructions from README:
   ```
   https://github.com/nethesis/dartagnan/blob/master/aramis/README.md
   ```

2. Copy the content of `dist` directory inside `/opt/dartagnan/www`

3. Configure [Auth0](https://auth0.com/) account and set `auth0-config.js` accordingly.