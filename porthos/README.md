Porthos
=======

Installation
------------

On CentOS 7 - porthos,

    yum -y --enablerepo=extras install epel-release
    yum -y install nginx php-fpm redis php-pecl-redis stunnel rsync
    systemctl enable nginx php-fpm redis@athos

On your local system,

    rsync -ai porthos/root/ root@porthos:/

On CentOS 7 - porthos,

    systemctl start nginx php-fpm redis@athos

YUM client
----------

HTTP repository metadata query:

    http://porthos.nethserver.org/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml

HTTP authentication is required to GET `repomd.xml`:

    http://<system_id>:<secret>@porthos.nethserver.org/<repo_version>/<repo_name>/<repo_arch>/repodata/repomd.xml


Redis DB format
---------------

The `repomd.php` script expects the following storage format in redis DB:

    key: <system_id>
    value: hash{ tier_id => <integer>, secret => <string> }

Troubleshooting
---------------

HTTP status codes for `repomd.xml`

* 404 - server_id not found
* 403 - server_id exists, but has disabled access or the given secret is bad
* 503 - redis connection failed, see `/var/log/nginx/porthos-php-error.log`
* 502 - php-fpm connection failed, see `/var/log/nginx/error.log`
* 500 - generic PHP error, see `/var/log/nginx/porthos-php-error.log`

TODO
----

- configure `yum-cron` automated daily updates
- configure `iptables`
- command `repo-pull-upstream`
- command `repo-shift-tier`
