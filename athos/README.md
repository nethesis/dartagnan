# athos

Athos implements a set of APIs to manage servers, subscriptions and alerting.
It's built on top of [Gorm](https://github.com/jinzhu/gorm) and [Gin-Gonic](https://github.com/gin-gonic)

## Build

Execute:
```
go get
go build
```

## Configuration

Example configuration:
```
{
    "db_host":"localhost",
    "db_port":"5432",
    "db_name":"dartagnan",
    "db_user": "dt_user",
    "db_password": "dt_password",
    "cors": {
        "origins": ["*"],
        "headers": ["Authorization", "Content-Type"],
        "methods": ["GET", "PUT", "POST", "DELETE"]
    },
    "auth0": {
        "domain": "dartagnan.auth0.com",
        "identifier_api":"http://my.nethserver.com:8080"
    },
    "paypal": {
        "client_id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "client_secret": "xxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    }
}
```

## Run

After build, execute:
```
./athos -c config.json
```
