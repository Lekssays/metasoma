syntax = "proto3";

package message;

option go_package = "github.com/Lekssays/ADeLe/autopeering/proto/peering";

message Ping {
    bool payload = 1;
}

message Pong {
    bool payload = 1;
}

message Request {
    string publickey = 1;
}

message Response {
    bool result = 1;
    string proof = 2;
    string signature = 3;
    string publickey = 4;
}
