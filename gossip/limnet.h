#include <stdbool.h>
#include <stdint.h>
#include <stddef.h>

#define ip_t uint32_t

#ifdef __cplusplus
extern "C" {
#endif

// All pointers are managed managed by the caller.
// The functions never retain those pointers,
// so they can be freely deallocated or reused as soon as the function call returns

// initializes the LiMNet subsystem
// activate_sniffer: true for real packet sniffing, false to load preprocessed packet features with on_packet_received()
// returns: 0 if no errors, > 0 if a fatal error occurred and the initialization failed
int initialize(bool activate_sniffer);

// requires the LiMNet subsystem to terminate and waits for said termination to complete
void terminate();

// called when memories are received from a peer
// ips: buffer of num_entries ips (each 4 bytes)
// memories: buffer of num_entries compressed memories (each as big as returned by compressed_memory_size())
// num_entries: number of entries in the buffers (not buffer size)
void on_memories_received(ip_t* ips, void* memories, size_t num_entries);

// called when packet features are received based on a preprocessed trace
// ip_src: ip address of the source node for this packet
// ip_src: ip address of the destination node for this packet
// features: buffer containing the packet features (as big as returned by packet_features_size())
void on_packet_received(ip_t ip_src, ip_t ip_dst, float* features);

// returns at most the requested amount of ip addresses and corresponding memories,
// with priority to those with memories that should be broadcast soon
// ips: output buffer, will be filled with the returned ips, must be able to contain up to num_entries ips (each 4 bytes)
// memories: output buffer, will be filled with the returned memories, must be able to contain up to num_entries memories (each as big as returned by compressed_memory_size())
// num_entries: maximum number of ips and memories to put in the buffers
// returns: actual number of ips put in the buffer (possibly less than num_entries)
size_t get_memories_to_share(ip_t* ips, void* memories, size_t num_entries);

// returns the ip address of a random peer to which memories should be broadcast
ip_t get_random_peer();

// returns the size in bytes of a compressed memory vector
size_t compressed_memory_size();

// returns the size in bytes of packet features from a preprocessed trace
size_t packet_features_size();

#ifdef __cplusplus
}
#endif