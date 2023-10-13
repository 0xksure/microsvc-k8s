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

# Architecture
Everything is orchestrated using kubernetes and exposed through a load balancer for local development. 

# TODO
[ ] make service 1 make useful calculations
[ ] stand up a db to store state. Use postgresql rather than scylla. 
[ ] build and expose a frontend server or just a static web server 
