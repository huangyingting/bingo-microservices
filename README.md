# Bingo - Cloud Native Showcase Application
## Introduction
Bingo is a short url application based on microservice's architecture. It helps you create and share branded links with custom domains at scale. 

Why create this short url application for cloud native showcase - generating a short url could be simple or complex enough at the same time.

Below high level architecture diagram shows each componnent
![high-level-design](./docs/images/Bingo-Design.svg)

- **BF** - Singal Page Application(SPA), a frontend application written in Javascript with React.js who provides web UI to create and manage short url.
- **BS** - Bingo service, core service written in Golang provides API to create, update, delete and get short url, this service also serves as HTTP server to redirct URL.
- **BE** - Bingo Extract service, is a service written in Python where you can find out where the shortened URL will take you to before clicking on the link. 
- **BI** - Bingo Intelligence service, is a service written in Golang that provides detailed statistics for your short url, such as clicks, end user geo distributions etc.
- **BG** - Bingo Geo service, geo location service written in Golang translates end user IP address to country and city.
- **GoWitness** - [A website screenshot utility](https://github.com/sensepost/gowitness) written in Golang.

## Features
- Bingo supports multiple database types, including sqlite(testing purpose), mysql, postgres, sql server and mongodb.

- REST, GRPC and Websocket. REST API is exposed for external use, internally, BE and BI services support GRPC. Websocket is enabled on BS service to send back server updated message.

- Message queue, rabbitmq provides message exchange between BS and BI services. BS service publishes click stream to rabbitmq, and BI service subscribes to click stream.

- Unique alias generator, when create a short url, a 7 characters alias will be generated, bingo supports a [sonyflake](https://github.com/sony/sonyflake) based algorithm to generate this alias among at most 16 machines.

- Distributed lock, when generate alias, the algorithm requires a unique machine ID so no duplicated alias will be genereated among each machine, ETCD distributed lock is used here to provide this capability.

- Caching and bloom filter support, on top of database layer, BS service also uses redis for caching and bloom filter to protect database from overloading.

- Tag suggestion, BS uses elasticsearch to store short url tags and provide tag suggestion based on existing tags - a seamless search as you type experience.

- Natural language processing, BE uses NLP to extract web URL's keywords and summary.

## How to build and run
### Build and run locally
1. git clone https://github.com/huangyingting/bingo-microservices
2. cd bingo-microservices
3. make docker
4. docker-compose up
5. Visit URL http://localhost:8080 to login