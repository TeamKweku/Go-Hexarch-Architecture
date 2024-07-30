# _Redesigning_ Code Odessey Using a Modular Monolithic Architecture Employing the Hexagonal Architecture Pattern, Type-Driven Development and Production-Ready Approach in Golang

> DISCLAIMER: Code in this repository is intended for educational purpose and a continuation of the [Code Odyssey](https://github.com/TeamKweku/code-odessey) project. Instead of refactoring the existing codebase, we are redesigning the project using a modular monolithic architecture employing the hexagonal architecture pattern, type-driven development and a production-ready approach just as a learning process for various concepts and best practices in software development. Also note the the goal of this project is to be able to use the ports and adpater nature to integrate with many technologies and services as possible.

This repository aims to showcase advanced techniques for building a production-ready application in Golang. Mostly inspired by the [realworld-go](https://github.com/AngusGMorrison/realworld-go) by `AngusGMorrison`. The project is primarily a learning exercise for various Golang and DevOPS concepts, I attempt to learn during the journey of building this project.

## Purpose of Application

In my experience aside building a ecommerce project that kind of involves incorporating alot of software development concepts, I found creating a blog app to be a good way to learn and practice various software development concepts. Furthermore the implemenation of this project aims to use the Hexagonal Architecture pattern hence seperating the business logic from the infrastructure logic. This will make it easy to swap out the infrastructure logic with another one without affecting the business logic.The example code uses `gRPC` and `gRPC gateway` to implement a `RESTFUL` API. The project also uses `PostgreSQL` as the database and `Redis` as the cache store. The project also uses `Docker` and `Docker Compose` to containerize the application and its dependencies. But since business logic won't be tied so much with the infrastructure logic, in due time we can swap out the infrastructure logic with technologies like `MongoDB`, `RabbitMQ`, `Kafka`, `AWS`, `GCP`, `Azure`, `OpenTelemetry`, `Elastic Search` etc.

## Hexagonal Architecture

This project is just an implementation of the `Hexagonal Architecture` but for an explanation of this architecture and its benefits, you can check out [realworld-go](https://github.com/AngusGMorrison/realworld-go), where he references some interest articles on the topic from `Netflix` and `Uber`. For a video explanation on the architecture checkout [How To Structure Your Go App - Full Course](https://www.youtube.com/watch?v=MpFog2kZsHk&list=PL7g1jYj15RUPjxpD_PDt8L7IlA-VpT0t8) playlist.

## How Project is Structured

> COMING SOON..

# Progress

Here's what's been implemented so far.

## Code Odessey

- [ ] Users
- [ ] Authentication
- [ ] Profiles
- [ ] Articles
- [ ] Comments
- [ ] Tags
- [ ] Favorites
- [ ] Reactions
- [ ] Search

## Productionization

- [ ] CI pipeline
- [ ] Optimized Docker image
- [ ] First-class error handling
- [ ] Configuration
- [ ] Linting
- [ ] Extensive, concurrent unit test suite
- [ ] Health checks
- [ ] Streamlined local development experience
- [ ] Optimistic concurrency control
- [ ] Concurrent integration tests
- [ ] Structured logging
- [ ] Metrics
- [ ] Tracing
- [ ] API documentation
- [ ] Local Deployment
