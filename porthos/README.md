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

If `tier_id` is not a number, `stable/` is served instead.

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

## Porthos repository autoupdate policy

Changes to upstream repostories are recorded as repository snapshots, according to the
configured cron table.

A snapshot directory contain the differences with the upstream repository at snapshot
creation time. It is created by the `rsync --backup` command. Differences may consist of
`repomd.xml` file and related YUM repository metadata files. If the upstream
repository changes include also some RPM deletions, the snapshot directory contains
also RPM files.

When a client accesses the `autoupdate/` contents it is possible to
decide what point in time it can see, on a client `tier_id` and `version` basis.

The **default** policy is to return the repository snapshot contents of the last Monday, or
the previous one, depending on the tier number.

Custom policies can be configured by adding items to the `$config['autoupdate_policy']` array.

Example 1:

```php
$config['autoupdate_policy'] = array(
    '7.6.1810/0' => 'head',
    '7.6.1810/*' => 'empty',
);
```

This setting says that clients of tier 0 requesting version 7.6.1810 can see the
`head` repository state. Clients of other tiers (identified by `*`) always see
an `empty` repository.

Example 2:

```php
$config['autoupdate_policy'] = array(
    '7.6.1810/*' => 'default',
    '7.6.1810/2' => 'fixed/d20191030',
);
```

In this case clients of tier 2 always see the "fixed" repositories state as
they were recorded by snapshots up to 2019-10-30.
Other clients see the repository state according to the default policy (note
that the corresponding line can be omitted because it already corresponds to the
default policy).

## Repository management commands

The `repo-*` are a set of Bash commands that include (source) the configuration
from `/etc/porthos/repos.conf`. Upstream YUM rsync URLs are defined there.

The following commands are executed automatically, as defined in `porthos.cron`:

- `repo-bulk-pull` creates a snapshot date-dir (e.g. `d20180301`) under
  `dest_dir` with differences from upstream repositories.
- `repo-bulk-cleanup` erases stale snapshots directories

The following commands are designed for Porthos initialization, to recover from errors, or implement low-level actions:

- `repo-bulk-hinit` runs initial synchronization from upstream repositories (-f disables the check for already existing directories)
- `repo-head-init`  synchronization of head from a specific upstream repo
- `repo-head-rollback` rolls back a repository to its latest snapshot
- `repo-snapshot-create` create a new repository snapshot
- `repo-snapshot-delete` delete repomd.xml from a given repository snapshot
- `repo-rpm-lookup`  seek the given RPM in every snapshot for a given repository
- `xrsync` run rsync safely, trying to repeat the operation if it fails

A **rollback action** for a given repository consists into reverting its most
recent snapshot state, by moving the snapshot YUM metadata and removed/changed
RPMs to the head position. For instance:

    repo-head-rollback 7.7.1908/nethserver-updates/x86_64 -n

The command above attempts to roll back the `head/` directory to the most
recent, non-empty snapshot of `nethserver-updates`. The command can be invoked
multiple times, but it fails as soon as no snapshot is found, or if an invalid
repository identifier is issued. Remove the `-n` flag (rsync dry run) to
actually change the files.

Some times it is desirable to re-sync the head repository, without generating a
new snapshot, like `repo-snapshot-create` does. That happens if an upstream repo
was fixed quickly and the bogus RPM never entered any snapshot. In that case run
`repo-head-init` as follow:

    repo-head-init -n -f 7.6.1810/nethserver-updates/x86_64

The `-n` flag preserves local files from deletion, whilst `-f` forces the
command to run even if the repository was already initialized.

If one or more snapshots contain a bogus RPM it is possible to delete the whole
repository metadata (repomd.xml) file with the following command:

    repo-snapshot-delete d20190702/7.6.1810/nethserver-updates/x86_64 d20190630/7.6.1810/nethserver-updates/x86_64

The correct snapshot name can be found starting from the RPM name with:

    repo-rpm-lookup bogus-rpm-1.2.3-1.ns7.noarch.rpm
    d20190702/7.6.1810/nethserver-updates/x86_64
    d20190630/7.6.1810/nethserver-updates/x86_64

The two commands can be combined together with `xargs`:

    repo-rpm-lookup bogus-rpm-1.2.3-1.ns7.noarch.rpm | xargs -- repo-snapshot-delete

If the RPM is found under `head/`, `repo-snapshot-delete` safely ignores it.

## Automated schedule

The management commands are executed at specific days of the week, as specified
by the ``/etc/cron.d/porthos.cron`` crontab.

## New minor release checklist

When upstream releases a new minor version, 

- fix the `repos.conf` configuration file with new release number / mirror location
- run the initial synchronization `repo-bulk-hinit`
- add the new release to `config-porthos.php`
