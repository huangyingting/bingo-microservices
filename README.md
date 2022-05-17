# Bingo - Cloud Native Showcase Application
## Introduction
Bingo is a short url application based on microservice's architecture. It helps you create and share branded links at scale. 

The purpose of creating this short url application for cloud native showcase is, the business logic of generating and using a short url could be simple or complicated at the both time, which gives us enough space to demo a few of cloud native design principles. 

Below high level architecture diagram shows all related componnents and services
![high-level-design](./docs/images/Bingo-Design.svg)

- **BF** - Singal Page Application(SPA), a frontend application written in Javascript with React.js who provides web UI to create and manage short url.
- **BS** - Bingo service, core service written in Golang provides API to create, update, delete and get short url, this service also serves as HTTP server to redirct URL.
- **BE** - Bingo Extract service, is a service written in Python where you can find out where the shortened URL will take you to before clicking on the link. 
- **BI** - Bingo Intelligence service, is a service written in Golang that provides detailed statistics for your short url, such as clicks, end user geo distributions etc.
- **BG** - Bingo Geo service, geo location service written in Golang translates end user IP address to country and city.
- **GoWitness** - [A website screenshot utility](https://github.com/sensepost/gowitness) written in Golang.

## Features
- Multiple database types support, including sqlite(testing purpose), mysql, postgres, sql server and mongodb.

- Observability support, all microservices support logging, metrics(prometheus exporter), and tracing (jaeger based distributed tracing support).

- AAD authentication & authorization support, BS APIs are protected by oauth so only valiated user can call those APIs to create, edit, update and delete.

- REST, GRPC and Websocket. REST API is exposed for external use, internally, BE and BI services support GRPC. Websocket is enabled on BS service to send back server updated message.

- Message queue, rabbitmq provides message exchange between BS and BI services. BS service publishes click stream to rabbitmq, and BI service subscribes to click stream.

- Unique alias generator, when create a short url, a 8 characters alias will be generated, bingo supports a [sonyflake](https://github.com/sony/sonyflake) based algorithm to generate this alias.

- Distributed lock, when generate alias, the algorithm requires a unique machine ID so no duplicated alias will be genereated from each machine, ETCD distributed lock is used here to provide this capability.

- Caching and bloom filter support, on top of database layer, BS service also uses redis for caching and bloom filter to protect database from overloading.

- Tag suggestion, BS uses elasticsearch to store short url tags and provide tag suggestion based on existing tags - a seamless search as you type experience.

- Natural language processing, BE uses NLP to extract web URL's keywords and summary.

## How to build and run
### Prerequisite
Bingo replies on AAD to provide authentication & authorization, the application itself is pre-configured with an AAD tenant already, if you prefer to use your own AAD tenant, below are the brief steps
- Create an AAD or AAD B2C tenant
- Associate a custom domain name
- Add a few of users into the tenant
- Register AAD client application(used by BF) and api application(used by BS)
- Create scopes for api application then assign those socpes as API permissions to client application
- Record client application id, scopes and oauth2 endpoint, those information are required to configure client app and BS service. For more details, refer to js/bf/src/Global.js and go/app/bs/configs/config.yaml

References
1. [Register a Microsoft Graph application](https://docs.microsoft.com/en-us/azure/active-directory-b2c/microsoft-graph-get-started?tabs=app-reg-ga)
2. AAD [Single-page application: App registration](https://docs.microsoft.com/en-us/azure/active-directory/develop/scenario-spa-app-registration), AAD B2C [Register a single-page application (SPA) in Azure Active Directory B2C](https://docs.microsoft.com/en-us/azure/active-directory-b2c/tutorial-register-spa)
3. AAD [Protected web API: App registration](https://docs.microsoft.com/en-us/azure/active-directory/develop/scenario-protected-web-api-app-registration), AAD B2C [Add a web API application to your Azure Active Directory B2C tenant](https://docs.microsoft.com/en-us/azure/active-directory-b2c/add-web-api-application?tabs=app-reg-ga)

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
4. Start all services 
    ```
    make up
    ```
5. Visit URL http://localhost:8080 to login with your AAD user account.