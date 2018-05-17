# Porthos

Porthos provides YUM repositories authentication, authorization and progressive 
release of updates to individual clients.

The basic rule of Porthos releases is that the **minimum age** of an RPM from
stable repositories is **one week**.

## Installation

On CentOS 7 - porthos,

    yum -y --enablerepo=extras install epel-release
    yum -y install nginx php-fpm redis php-pecl-redis stunnel rsync
    systemctl enable nginx php-fpm redis@athos

On your local system,

    rsync -Cai --no-super --no-o --no-g porthos/root/ root@porthos:/

On CentOS 7 - porthos, configure and apply SELinux context

    systemctl daemon-reload
    semanage fcontext -a -t httpd_sys_content_t '/srv/porthos(/.*)?'
    restorecon -vvRF /srv/porthos

Allow SELinux access to redis Unix socket from php-fpm

    setsebool -P daemons_enable_cluster_mode 1

Run initial synchronization, then start daemons

    repo-bulk-hinit
    systemctl start nginx php-fpm redis@athos

## Certbot configuration

    yum install certbot
    systemctl enable --now certbot-renew.timer
    certbot certonly -n --cert-name porthos --allow-subset-of-names \
        --agree-tos --email nethinfra@nethesis.it \
        --webroot -w /srv/porthos/certbot/ \
        -d $(hostname) -d mirrorlist.nethserver.com
    restorecon -vvRF /etc
    
Edit `/etc/nginx/conf.d/porthos.conf` and `/etc/nginx/conf.d/mirrorlist.conf`
and uncomment the following line:

    # include porthos-certbot.conf;

Restart nginx:

    systemctl restart nginx

## YUM client

Mirrorlist query format

    http://mirrorlist.nethserver.com/?repo=nethserver-base&arch=x86_64&nsversion=7.4.1708&usetier=yes

The `mirrorlist.nethserver.com` virtual host returns the list of available YUM
repository mirrors for the given parameters.

- `nsversion` full version number `X.Y.Z` 
- `repo` repository name
- `usetier` whether tier repository is desired or not. If present, `$YUM0`,
  `no`, `0` and empty string are mapped to `false`, any other value is mapped to
  `true`
- `arch` system architecture, like `x86_64`

HTTP repository metadata query (HTTP Basic authentication required):

    https://m1.nethserver.com/autoupdate/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml
    https://m1.nethserver.com/stable/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml

* `autoupdate` returns data from the tier associated to the credentials provided
* `stable` always returns data from `t0`

## HTTP status codes

* 401 - authorization required
* 404 - resource not found
* 403 - bad credentials or `tier_id` is false (disabled)
* 503 - redis connection failed, see `/var/log/nginx/porthos-php-error.log`
* 502 - php-fpm connection failed, see `/var/log/nginx/error.log`
* 500 - generic PHP error, see `/var/log/nginx/porthos-php-error.log`

## Redis DB format

The `repomd.php` script expects the following storage format in redis DB:

    key: <system_id>
    value: hash{ tier_id => <integer>, secret => <string> }

If `tier_id` is not set, the access is denied (403 - forbidden). For instance to create a key on athos

    redis-cli -p PORT
    redis-cli PORT> HMSET 0ILD29RH-D78A-C444-1F82-EE92-3211-FC47-43AD-DQFD tier_id 2 secret S3Cr3t

## Repository management commands

The `repo-*` are a set of Bash commands that include (source) the configuration
from `/etc/porthos.conf`. Upstream YUM rsync URLs are defined there.

- `repo-bulk-hinit` runs initial synchronization from upstream repositories (-f disables the check for already existing directories)
- `repo-bulk-pull` creates a snapshot date-dir (e.g. `d20180301`) under
  `dest_dir` with differences from upstream repositories. Set `t0` to point at
  it.
- `repo-bulk-shift [N]` updates `t1` ... `tN` links by shifting tiers of one position
  the optional `N`
- `repo-bulk-cleanup` erases stale tier snapshots

The following commands should not be invoked directly. They are intended to be
called by the commands above.

- `repo-head-init`  initial synchronization from a specific upstream repo
- `repo-tier-pull`  upstream snapshot for a specific repo
- `xrsync` run rsync and try to repeat the operation if fails

## Automated schedule

The management commands are executed at specific days of the week, as specified
by the ``/etc/cron.d/porthos.cron`` crontab.

## New minor release checklist

When upstream releases a new minor version, 

- fix the `repos.conf` configuration file with new release number / mirror location
- run the initial synchronization `repo-bulk-hinit`