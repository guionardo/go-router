# go-router
Friendly router for REST APIs

This is a router library to create an descriptive API, that abstract the http layer and provides simpler service/use case/business rules.

## Example:

```go

r:=router.New(WithTitle("Example API"),WithVersin("0.0.1"))
r.Get("/ping",)
```

## Components of the router

### Router

Base of the router 

## Links

* [Best Go web frameworks for 2025](https://blog.logrocket.com/top-go-frameworks-2025/)


## RequestData

### Path values

Using path templates, the fields will be read from the request path.

* /entity/:id with a request /entity/2 will inject {"id":"2"} 

### Headers
