# Kubernetes microservices Demo app

The app sets out to demonstrate how to run microservices in a kubernetes cluster.  

## Setup

### Install just
> https://github.com/casey/just

### Install migrate 
https://github.com/golang-migrate/migrate/tree/master/cmd/migrate


# Services 

## Service 1 (to be renamed)
Service 1 is a go service that simply exposes an RPC server. The server is only exposed within the cluster and is not accessible from the outside. 

## Service 2 (to be renamed)
Service 2 is an exposed service that leverages service 1 through RPC in order to execute smaller tasks

## Github app (golang)
This is a github app that listens to comments on issues and pull requests. It is activated when the pattern
```
$<string>:<int>$
```
is observed in a create action. It will then parse the string and int which in this context is a denomination and an amount. This is basically the bounty for the given issue or pull request.

If a new bounty is observed it prints a redirect url to  the frontend service that handles signing of the bounty creation. If successful the user will be redirected to to the issue or pull request. 

## Frontend (sveltekit)
The frontend service is a ssr sveltekit application running in v8. It is mainly responsible for 
1. Signing a new bounty with a browser wallet
2. Taking arguments to link identities

After 1 it records the signed to the chosen message bus and redirects the user. 

After successfully linking an identity it sends the message to the chose message bus. The event is picked up by the identity service. 

## Identity Service (not yet implemented)


## Identity app ()
It is resposible for linking identities and in general authorization. It will store the links in it's own database.

# Architecture
A microservice architecture is chosen because it reduces the logic whitin the services. Kubernetes is being leveraged as the infrastructure to orchestrate the containers and the communication between them.

# Technologies
- Kubernetes
- Postgresql 
- Golang
- Node (sveltekit)
- Docker 

## Potential technology
- Rust for identity service
- Kafka (wo zookeeper) as a "message bus"
- A database for heavy read access 

# TODO
- [x] make service 1 make useful calculations
- [x] stand up a db to store state. Use postgresql rather than scylla. 
- [x] build and expose a frontend server or just a static web server 
- [ ] Stand up message bus 
- [ ] Specify ingress and egress rules for k8s clsuter 
