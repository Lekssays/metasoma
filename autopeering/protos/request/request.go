syntax = "proto3";

package message;

option go_package = "github.com/Lekssays/ADeLe/autopeering/proto/request";

message Ping {

}

message Message {
    Memory mergedMemory = 1;
    repeated Memory parents = 2;
    repeated string peeringProofs = 3;
}