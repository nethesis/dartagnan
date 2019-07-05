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

HTTP repository metadata query (SSL + HTTP Basic authentication required):

    https://m1.nethserver.com/autoupdate/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml
    https://m1.nethserver.com/stable/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml

HTTP repository metadata query (username as path token, SSL not required):

    http://m1.nethserver.com/autoupdate/<username>/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml
    http://m1.nethserver.com/stable/<username>/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml


* `autoupdate` returns data from the tier associated to the credentials provided
* `stable` always returns data from `t0`

## HTTP status codes

* 401 - authorization required
* 404 - resource not found
* 403 - bad credentials or access denied conditions
* 503 - redis connection failed, see `/var/log/nginx/porthos-php-error.log`
* 502 - php-fpm connection failed, see `/var/log/nginx/error.log`
* 500 - generic PHP error, see `/var/log/nginx/porthos-php-error.log`

## Porthos repository access control

Access permissions to Porthos repositories are checked by `auth.php`.

A special `subscription.json` API endpoint used by the Software center clients
is implemented by `subscription.php`.

All HTTP requests to access YUM repositories must be authenticated with HTTP
Basic Auth. If the special `$config['legacy_auth']` is enabled, the HTTP
username is considered a valid access token, and can be set as password value too.

The access token can be passed also as an URL path component, like:

    http://m1.nethserver.com/stable/<username>/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml

The `auth.php` and `subscription.php` scripts expects the following record
format in redis DB:

    key: <system_id>
    value: hash{ tier_id => <value>, secret => <value>, icat => <value> }

For instance to create a key on athos

    redis-cli -p PORT
    redis-cli PORT> HMSET 0ILD29RH-D78A-C444-1F82-EE92-3211-FC47-43AD-DQFD tier_id 2 secret S3Cr3t icat cat1,cat2,cat3

### `tier_id` field

The `tier_id` value should be a number. If the value is negative, the tier
number is calculated by an hash function, based on the system identifier. 

If `tier_id` is not a number, both `auth.php` and `subscription.php` reply with
403 - forbidden.

### `icat` field

The `icat` field is a string of a comma separated list of YUM category
identifiers (refer to the repository comps/groups for valid names). Its purpose
is to show the products entitlement on the Software center page. It is used by
`subscription.php` to return the included/excluded YUM categories list to the
client. See also the `$config['categories']` parameter to configure it. This
field is ignored by `auth.php`.

If `icat` field is not set, the `subscription.php` replies with 403 - forbidden.

### `secret` field

If `secret` field is not set, both `auth.php` and `subscription.php` reply with
403 - forbidden, unless `$config['legacy_auth']` is enabled.


## Repository management commands

The `repo-*` are a set of Bash commands that include (source) the configuration
from `/etc/porthos/repos.conf`. Upstream YUM rsync URLs are defined there.

The following commands are executed automatically, as defined in `porthos.cron`:

- `repo-bulk-pull` creates a snapshot date-dir (e.g. `d20180301`) under
  `dest_dir` with differences from upstream repositories. It sets `t0` to point at
  it.
- `repo-bulk-shift [N]` updates `t1` ... `tN` links by shifting tiers of one position
  the optional `N` parameter creates missing links up to N - 1.
- `repo-bulk-cleanup` erases stale snapshots directories

The following commands are designed for Porthos initialization, to recover from errors, or implement low-level actions:

- `repo-bulk-hinit` runs initial synchronization from upstream repositories (-f disables the check for already existing directories)
- `repo-head-init`  initial/override synchronization of head from a specific upstream repo
- `repo-head-rollback` roll back head to a previous snapshot for a specific repo
- `repo-tier-pull`  create a new upstream snapshot for a specific repo
- `repo-tier-delete`  delete repomd.xml from a given tier or snapshot
- `repo-rpm-lookup`  seek the given RPM in every snapshot for a given repository
- `xrsync` run rsync safely, trying to repeat the operation if it fails

A **rollback action** for a given repository consists into seeking the
most recent snapshot and moving it back to the head position. YUM metadata and
removed RPMs are merged, reverting the head to the past snapshot state. For
instance:

    repo-head-rollback 7.6.1810/nethserver-updates/x86_64

The command above rolls back the `head/` directory to the most recent, non-empty
snapshot of `nethserver-updates`. The command can be invoked multiple times, but
it fails as soon as no snapshot is found, or if an invalid repository identifier
is issued.

Some times it is desirable to re-sync the head repository, without generating a
new snapshot, like `repo-tier-pull` does. That happens if an upstream repo was
fixed before being shifted. In that case run `repo-head-init` as follow:

    repo-head-init -n -f 7.6.1810/nethserver-updates/x86_64

The `-n` flag preserves local files from deletion, whilst `-f` forces the
command to run even if the repository was already initialized.

If one or more snapshots contain a bogus RPM it is possible to delete the whole
repository metadata (repomd.xml) file with the following command:

    repo-tier-delete d20190702/7.6.1810/nethserver-updates/x86_64 d20190630/7.6.1810/nethserver-updates/x86_64

The correct snapshot (or tier) name can be found starting from the RPM name with:

    repo-rpm-lookup bogus-rpm-1.2.3-1.ns7.noarch.rpm
    d20190702/7.6.1810/nethserver-updates/x86_64
    d20190630/7.6.1810/nethserver-updates/x86_64

The two commands can be combined together with `xargs`:

    repo-rpm-lookup bogus-rpm-1.2.3-1.ns7.noarch.rpm | xargs -- repo-tier-delete

If the RPM is found under `head/`, `repo-tier-delete` safely ignores it.

## Automated schedule

The management commands are executed at specific days of the week, as specified
by the ``/etc/cron.d/porthos.cron`` crontab.

## New minor release checklist

When upstream releases a new minor version, 

- fix the `repos.conf` configuration file with new release number / mirror location
- run the initial synchronization `repo-bulk-hinit`
- add the new release to `config-porthos.php`
