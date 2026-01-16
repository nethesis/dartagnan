# Local Deployment

## Configuration File Changes

### **aramis/config/index.js**

- add the proxy to simulate the behavior performed in production environment by *nginx*

    ```
    ...
    proxyTable: {
        '/api': {
            target: 'http://localhost:8080',
            changeOrigin: true
        }
    }
    ...
    ```
- change the startup port of the web interface to avoid conflicts with the container that exposes backend services

    ```
    ...
    port: 5000,
    ...
    ```

### **aramis/static/config/config.js**

- set *client_id*, *domain* and *audience* of the application created and configured on [Auth0](https://auth0.com/). You will need to define an API for the audience and a new Single Page Application to handle the domain.

    ```
    ...
      AUTH0_CLIENT_ID: '<YOUR-APPLICATION-CLIENT-ID>',
      AUTH0_DOMAIN: '<YOUR-APPLICATION-DOMAIN>',
      AUTH0_CALLBACK: window.location.origin+'/callback',
      AUTH0_AUDIENCE: '<YOUR-API-AUDIENCE>',
    ...
    ```

N.B. To speed up the authentication mechanism, also define a user with username and password in the *User Management > Users* section.

- define a new application on the [PayPal Developer Portal](https://developer.paypal.com/home/) in the *App and Credentials > Create App* section, from which you copy the *client-id*. Also verify the presence of two account types (*personal* and *business*) in the *Testing Tools > Sandbox Accounts* section, which will be used to simulate the merchant account (where the purchase is made) and a generic personal account (with which to make the purchase using a credit that can be modified at will once exhausted).

    ```
    ...
      PAYPAL_SANDBOX: '<YOUR-PAYPAL-APPLICATION-CLIENT-ID>',
      PAYPAL_PRODUCTION: '',
    ...

    ```

### aramis/src/components/directives/RenewButton.vue

- to avoid issues with the PayPal library used, you need to disable the call that is executed for creating web profiles (which allows you to customize the purchase modal shown by PayPal with your logo, information, ...). To do this you need to modify the payload sent in the payment creation call

    ```
    ...
            payment: function (data, actions) {
              return actions.payment.create({
                payment: {
                  intent: "sale",
                  payer: {
                    payment_method: "paypal"
                  },
                  transactions: [
                    {
                      amount: {
                        total:
                          Math.round(
                            (context.currentPlan.price +
                              (context.currentPlan.price *
                                context.billingInfo.tax) /
                                100) *
                              100
                          ) / 100,
                        currency: "EUR",
                        details: {
                          subtotal:
                            Math.round(context.currentPlan.price * 100) / 100,
                          tax:
                            Math.round(
                              ((context.currentPlan.price *
                                context.billingInfo.tax) /
                                100) *
                                100
                            ) / 100,
                        },
                      },
                      item_list: {
                        items: [
                          {
                            name: context.currentPlan.code,
                            description: context.currentPlan.name,
                            sku: context.obj.uuid,
                            price: Math.round(context.currentPlan.price * 100) / 100,
                            currency: "EUR",
                            quantity: "1",
                          },
                        ],
                      },
                    },
                  ],
                }
              });
            },
            ...
    ```

### deploy/roles/athos/files/database.sql

- disable user and database creation as it is already provided directly in the *docker-compose* when initializing the Postgres database container. The initialization script will be executed and the tables will be populated.

    ```
    -- CREATE database dartagnan;
    -- CREATE USER dtuser WITH PASSWORD 'SECRET';
    ALTER USER dtuser WITH SUPERUSER;
    GRANT ALL PRIVILEGES ON DATABASE dartagnan to dtuser;
    \connect dartagnan
    ...
    ```

### deploy/roles/athos/files/config.json

- set the correct country

    ```
    ...
        "billing": {
            "country": "Italy"
        },
    ...
    ```

- Postgres database configuration

    ```
    ...
        "database": {
            "host": "host.containers.internal",
            "port": "5432",
            "name": "dartagnan",
            "user": "dtuser",
            "password": "dtpassword"
        },
    ...
    ```

- Redis cache configuration

    ```
    ...
      "redis": {
            "host": "host.containers.internal",
            "port": "16379"
        }, 
    ...
    ```

- CORS policy configuration

    ```
    ...
        "origins": [
            "http://localhost:5000"
        ],
    ...
    ```

- Auth0 authentication configuration

    ```
    ...
        "auth0": {
            "domain": "<YOUR-APPLICATION-DOMAIN>",
            "audience": "<YOUR-API-AUDIENCE>"
        },
    ...
    ```

- PayPal configuration

    ```
    ...
        "paypal": {
            "sandbox": true,
            "client_id": "YOUR-PAYPAL-APPLICATION-CLIENT-ID",
            "client_secret": "YOUR-PAYPAL-APPLICATION-SECRET-ID"
        },
    ...
    ```


## Startup Procedure

1. **Copy necessary files**: copy the files contained in this folder (*podman-compose.yml* and *Containerfile*) to the main folder of the repository
    ```
    cp ./podman-compose.yml ./../..
    cp ./Containerfile ./../..
    ``` 
2. **Start the frontend**: go to the ```./aramis``` folder and run the two commands ```npm install``` and ```npm run dev```
3. **Start the backend**: run from the main repository folder the command ```podman-compose up -d```, which will build and start the containers for the *Athos* backend and related services (Redis cache and Postgres database)

### ⚠️ IMPORTANT - Remove Local Deployment Changes Before Committing

Before committing and pushing your new developments to the repository, **make sure to remove all modifications made exclusively for local deployment**:

- **Revert changes to configuration files** (Auth0 credentials, PayPal sandbox keys, database credentials, CORS origins, etc.)
- **Restore original settings** in:
  - `aramis/config/index.js` (proxy configuration)
  - `aramis/static/config/config.js` (Auth0 and PayPal settings)
  - `deploy/roles/athos/files/config.json` (database, Redis, CORS, Auth0, PayPal settings)
  - `deploy/roles/athos/files/database.sql` (database initialization)

Failure to do so may expose sensitive credentials or break the production build process.

## Production Environment Deployment Test

After developing a new feature or a fix locally, you can verify its compatibility with the production environment. To do this, follow these steps:

1. copy the *Containerfile* file contained in this folder to the main folder of the repository
    ```
    cp ./Containerfile ./../..
    ```

2. go to the main folder of the project and run the command
    ```
    podman build --target build -t dartagnan_athos_build .
    ```
    This command will build the image containing only the operating system present in the production environment and the copy of the source files, without actually executing the Athos backend.

3. once the image has been built, run the command
    ```
    podman run -v $(pwd):/build:Z dartagnan_athos_build sh -c "cd /build/athos && go get && go build"
    ```
    This command executes the compilation of the Athos backend on a container that matches the operating system present on the production machine.

4. once the previous command has been executed (some warnings may appear, but the important thing is that the command is actually executed and terminates correctly), to verify the successful outcome you will need to check the presence of the binary executable ```athos``` inside the folder ```/athos```. Run the command:
    ```
    ls -l ./athos | grep "athos"
    ```
    and verify the presence of the ```athos``` executable.