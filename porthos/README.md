Porthos
=======

Installation
------------

On CentOS 7 - porthos,

    yum --enablerepo=extras install epel-release
    yum install nginx php-fpm redis php-pecl-redis stunnel
    systemctl enable nginx php-fpm redis@athos stunnel@athos

On your local system,

    rsync -ai porthos/root/ root@porthos:/

On CentOS 7 - porthos,

    systemctl start nginx php-fpm redis@athos stunnel@athos

YUM client
----------

HTTP repository metadata query:

    http://porthos.nethserver.org/<server_id>/<repo_version>/<repo_name>/<repo_arch>/repodata/repomod.xml
    
Example:

    http://porthos.nethserver.org/0ILD29RH-D78A-C444-1F82-EE92-3211-FC47-43AD-DQFD/7.4.1708/nethserver-updates/x86_64/repodata/reppomd.xml

Troubleshooting
---------------

HTTP status codes

* 404 - server_id not found
* 403 - server_id exists, but has disabled access
* 503 - redis connection failed, see `/var/log/nginx/porthos-php-error.log`
* 502 - php-fpm connection failed, see `/var/log/nginx/error.log`
* 500 - generic PHP error, see `/var/log/nginx/porthos-php-error.log`

TODO
----

- configure `yum-cron` automated daily updates
- configure `iptables`
- command `repo-pull-upstream`
- command `repo-shift-tier`
