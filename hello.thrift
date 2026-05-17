namespace go hello
namespace py generated

struct HelloRequest {
    1: string name,
}

struct HelloResponse {
    1: string msg,
}

service Hello {
    HelloResponse hello(1: HelloRequest request),
}