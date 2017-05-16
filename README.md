# IoT Monitor Sample
This is a sample application designed to illustrate the use of Go to build a cloud native microservice and to use Redis as a backing store cache.

## Features
The following features are demonstrated by this sample application:
* A single, transport-agnostic service provides functionality to both **gRPC** and **HTTP** endpoints.
* Extensive use of Go Kit and the _middleware_ pattern
* Use of **Prometheus** and Go Kit metrics to expose advanced analytics
* Use of Redis as a cache to store the most recent telemetry, status update, and device registrations from sample IoT devices.
* Use of sub-packages for the _server_ and _client_ applications.
* Use of protocol buffers code generation from a _.proto_ file.