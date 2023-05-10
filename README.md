# LoadBalancer
A Low Level Design of Load Balancer in Golang

# Functional Requirements
1. API load balancer 
2. Rate limit - good to have 
3. Support Stateless API 
4. Dynamic config 


# How to run the code 
$GOROOT/bin/go run main.go 

A MakeFile is provided which can be used to build the binary as well as run the unit tests
For building binary, run "make" from $GOPATH/loadbalancer
For running tests, run "make test" from $GOPATH/loadbalancer

A Docker image is also provided for building the image and running this in a container
