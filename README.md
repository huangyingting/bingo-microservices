# Bingo - Cloud Native Showcase Application
## Introduction
Bingo is an URL shortener application based on microservices architecture. It is designed to create and share tiny URLs at a scale.

The purpose of creating this URL shortener application for cloud native showcase is, the business logic of generating and using a short URL could be simple or complicated at both time, it gives us enough space to demo a few of cloud native design principles.

Below high-level architecture diagram shows all related components and services  
![high-level-design](./docs/images/Bingo-Design.svg)

- **BF** - Single page application (SPA), a frontend application written in react.js who provides modern web UI to create and manage short URLs.
- **BS** - Bingo service, core service written in golang provides API to create, update, delete and get short URLs, it is also a HTTP server to redirect URL, protect bot attack and generate summary for URL.
- **BE** - Bingo extract service, URL summary service written in python, where you can get an overview of the short URL before click it, it is useful in blocking malicious short URL.
- **BI** - Bingo intelligence service, is a service written in golang that provides detailed statistics for your short URL, such as clicks, end user geo distributions etc.
- **BG** - Bingo geo service, is a geo location service written in golang who translates end user IP address to country and city.
- **GoWitness** - [A website screenshot utility](https://github.com/sensepost/gowitness) written in Golang.

## Features
- Support multiple databases as the backend storage, including sqlite (testing purpose), mysql, postgres, microsoft sql server and mongodb.

- Observability support, all microservices support logging, metrics (prometheus exporter), and distributed tracing (jaeger based).

- AAD authentication & authorization support, BS APIs are protected by oauth so only validated user can call API.

- REST and grpc API as well as websocket support. REST API is exposed for external use, internally, BE and BI services support GRPC. Websocket is enabled on BS service to send back server updated message.

- Message queue, RabbitMQ provides message exchange between BS and BI services. BS service publishes clickstream to RabbitMQ, and BI service subscribes to clickstream data.

- Unique alias generator, when create a short URL, a 8 characters alias will be generated, Bingo supports a [sonyflake](https://github.com/sony/sonyflake) based algorithm to generate this alias.

- Distributed lock, when generate alias, the algorithm requires a unique machine ID so no duplicated alias will be generated from each machine, ETCD distributed lock is used here to provide this capability.

- Caching and bloom filter support, on top of database layer, BS service also uses redis for caching and bloom filter to protect database from overloading.

- Tag suggestion, BS uses elastic search to store short URL tags and provide tag suggestion based on existing tags - a seamless search as you type experience.

- Natural language processing, BE uses NLP to extract web URL's keywords and summary.

## How to build and run
### Prerequisite - AAD/AAD B2C
Bingo replies on AAD to provide authentication & authorization, follow below steps to register bingo frontend and backend app from AAD or AAD B2c tenant
- Create an AAD or AAD B2C tenant
- Register a SPA AAD client application(used by BF, frontend) and an API application(used by BS, backend)
- Create scopes for API application then assign those socpes as API permissions to SPA application
- Record client application id, scopes and oauth2 endpoint, those information are required to configure bf frontend and BS service.
- Add a few of users into the tenant

References
1. [Register a Microsoft Graph application](https://docs.microsoft.com/en-us/azure/active-directory-b2c/microsoft-graph-get-started?tabs=app-reg-ga)
2. AAD [Single-page application: App registration](https://docs.microsoft.com/en-us/azure/active-directory/develop/scenario-spa-app-registration), AAD B2C [Register a single-page application (SPA) in Azure Active Directory B2C](https://docs.microsoft.com/en-us/azure/active-directory-b2c/tutorial-register-spa)
3. AAD [Protected web API: App registration](https://docs.microsoft.com/en-us/azure/active-directory/develop/scenario-protected-web-api-app-registration), AAD B2C [Add a web API application to your Azure Active Directory B2C tenant](https://docs.microsoft.com/en-us/azure/active-directory-b2c/add-web-api-application?tabs=app-reg-ga)

### Prerequisite - Google reCAPTCHA
Bingo relies on [Google reCAPTCHA](https://www.google.com/recaptcha/about/) for bot protection, please create a reCAPTCHA from reCAPTCHA admin portal, and record SITE KEY as well as SECRET KEY.

### Build and run locally
1. Clone the repo, 
    ```
    git clone https://github.com/huangyingting/bingo-microservices
    ```
2. Switch directory to repo, 
    ```
    cd bingo-microservices
    ```
3. Build container images for each service
    ```
    make docker
    ```
4. Create a file named .env.local under js/bf directory with below data
    ```
    BF_SCOPES_PREFIX=APPLICATION_ID_URI/
    BF_CLIENT_ID=SPA_APPLICATION_CLIENT_ID
    BF_AUTHORITY=https://login.microsoftonline.com/AAD_TENANT_ID
    ```
5. Create a file named .env.local under go/app/bs directory with below data
    ```
    BS_RECAPTCHA_SITE_KEY=GOOGLE_RECAPTCHA_SITE_KEY
    BS_RECAPTCHA_SECRET_KEY=GOOGLE_RECAPTCHA_SECRET_KEY
    # AAD v1 https://sts.windows.net/AAD_TENANT_ID/
    # AAD v2 https://login.microsoftonline.com/AAD_TENANT_ID/v2.0
    BS_JWT_ISSUER=https://sts.windows.net/AAD_TENANT_ID/
    # AAD v1 api://APPLICATION_ID_URI
    # AAD v2 APPLICATION_ID_URI
    BS_JWT_AUDIENCE=APPLICATION_ID_URI
    BS_JWT_TID=AAD_TENANT_ID
    ```
6. Start all services 
    ```
    make up
    ```
7. Visit URL http://localhost:8080 to login with your AAD/AAD B2C user account.

### Github CI/CD
Bingo is integrated with Github Actions, four workflows are included
- Build & push docker images to github container registry
- CodeQL to scan and discover vulnerabilities across go, javascript and python scripts
- Prune untagged images
- Prune pull request images
For more details, check .github/workflows directory

### Kubernetes deployment
Bingo supports kubernetes, a full set of deployment includes postgresql, redis, rabbitmq and etcd (helm charts from bitnami), elasticsearch (operator from elasticsearch), nginx ingress controller, cert-manager as well prometheus (prometheus-community/prometheus).

Deployment scripts files are included in deploy/ folder, deployment/bingo folder has bingo related yaml files, there is a .env file needs to be customized to include below configurations

```
BS_RECAPTCHA_SITE_KEY=GOOGLE_RECAPTCHA_SITE_KEY
BS_RECAPTCHA_SECRET_KEY=GOOGLE_RECAPTCHA_SECRET_KEY
# AAD v1 https://sts.windows.net/AAD_TENANT_ID/
# AAD v2 https://login.microsoftonline.com/AAD_TENANT_ID/v2.0
BS_JWT_ISSUER=https://sts.windows.net/AAD_TENANT_ID/
# AAD v1 api://APPLICATION_ID_URI
# AAD v2 APPLICATION_ID_URI
BS_JWT_AUDIENCE=APPLICATION_ID_URI
BS_JWT_TID=AAD_TENANT_ID
BF_SCOPES_PREFIX=APPLICATION_ID_URI/
BF_CLIENT_ID=SPA_APPLICATION_CLIENT_ID
BF_AUTHORITY=https://login.microsoftonline.com/AAD_TENANT_ID
```
Next, run `bash deploy.sh` to deploy all bingo microservices.