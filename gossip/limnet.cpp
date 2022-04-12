#include <chrono>
#include <memory>
#include <mutex>
#include <random>
#include <unordered_map>
#include <variant>
#include <vector>
#include <concurrentqueue/blockingconcurrentqueue.h>
#include <pcapplusplus/PcapLiveDevice.h>
#include <torch/script.h>
#include <torch/torch.h>
#include "limnet.h"
#include "generated/training_info.h"

#define COMPRESSED_DTYPE c10::kHalf

struct Memory {
    std::mutex mutex;
    float score;
    torch::Tensor memory;
    uint8_t flags;

    Memory() {
        score = 0;
        memory = torch::zeros({(int64_t)MEMORY_DIM}, c10::TensorOptions().dtype(c10::kFloat));
        flags = 0;
    }

    Memory(torch::Tensor m) {
        score = 0;
        memory = m.clone();
        flags = 0;
    }
};
using memory_t = std::shared_ptr<Memory>;

struct packet_t {
    ip_t src;
    ip_t dst;
    struct {
        float src_is_unicast;
        float src_is_private;
        float dst_is_unicast;
        float dst_is_private;
        float packet_len;
        float protocols[NUM_PROTOS];
    } features;
};
struct shared_mem_t {
    std::vector<ip_t> ips;
    torch::Tensor memories;
};
using message_t = std::variant<std::unique_ptr<packet_t>, shared_mem_t>;

moodycamel::BlockingConcurrentQueue<message_t> message_queue;
std::unordered_map<ip_t, memory_t> memories;
std::mutex memories_mutex;

torch::jit::script::Module forgery_model;
torch::jit::script::Module packet_model;
torch::jit::script::Module merge_model;
torch::jit::script::Module packet_cls_model;
torch::jit::script::Module device_cls_model;

std::thread limnet_thread;
volatile bool needs_to_terminate = false;

std::default_random_engine rng;

void limnet_thread_func() {
    torch::NoGradGuard no_grad;
    while (!needs_to_terminate) {
        message_t message;
        if !(message_queue.wait_dequeue_timed(message, std::chrono::milliseconds(100))) {
            continue;
        }
        switch (message.index()) {
            case 0: {
                auto packet = std::get<0>(std::move(message));
                memory_t src_mem, dst_mem;
                if (!memories.contains(packet->src)) {
                    std::lock_guard lock(memories_mutex);
                    src_mem = memories[packet->src] = memory_t(new Memory());
                } else {
                    src_mem = memories[packet->src];
                }
                if (!memories.contains(packet->dst)) {
                    std::lock_guard lock(memories_mutex);
                    dst_mem = memories[packet->dst] = memory_t(new Memory());
                } else {
                    dst_mem = memories[packet->dst];
                }
                auto pkt_features = torch::from_blob(&packet->features, {(int64_t)sizeof(packet_t::features)}, c10::TensorOptions().dtype(c10::kFloat));
                auto new_memories = packet_model.forward({src_mem->memory, dst_mem->memory, pkt_features}).toTensorVector();
                auto packet_class = packet_cls_model.forward({new_memories[0], new_memories[1], pkt_features}).toInt();
                auto src_class = device_cls_model.forward({new_memories[0]}).toInt();
                auto dst_class = device_cls_model.forward({new_memories[1]}).toInt();
                {
                    std::lock_guard lock(src_mem->mutex);
                    src_mem->memory = new_memories[0];
                    src_mem->flags = src_class;
                    src_mem->score += 1;
                }
                {
                    std::lock_guard lock(dst_mem->mutex);
                    dst_mem->memory = new_memories[1];
                    dst_mem->flags = dst_class;
                    dst_mem->score += 1;
                }
                break;
            }
            case 1: {
                auto shared_mems = std::get<1>(message);
                auto num_entries = shared_mems.ips.size();
                std::vector<std::pair<size_t, memory_t> > updated_memories(num_entries);
                int64_t cutoff = 0;
                size_t pos = num_entries - 1;
                for (size_t i = 0; i < num_entries; i++) {
                    auto ip = shared_mems.ips[i];
                    if (memories.contains(ip)) {
                        updated_memories[cutoff++] = {i, memories[ip]};
                    } else {
                        auto mem = memory_t(new Memory(shared_mems.memories[i]));
                        mem->flags = device_cls_model.forward({mem->memory}).toInt();
                        mem->score += 1;
                        updated_memories[pos--] = {ip, mem};
                    }
                }
                if (cutoff > 0) {
                    auto indices = torch::empty({cutoff}, c10::TensorOptions().dtype(c10::kLong));
                    auto old_mems = torch::empty({cutoff, (int64_t)MEMORY_DIM}, c10::TensorOptions().dtype(c10::kFloat));
                    for (size_t i = 0; i < cutoff; i++) {
                        indices[i] = (int64_t)updated_memories[i].first;
                        old_mems[i] = updated_memories[i].second->memory;
                    }
                    auto new_mems = merge_model.forward({old_mems, shared_mems.memories, indices}).toTensor();
                    for (size_t i = 0; i < cutoff; i++) {
                        auto [index, mem] = updated_memories[i];
                        auto new_class = device_cls_model.forward({new_mems[index]}).toInt();
                        std::lock_guard(mem->mutex);
                        mem->memory = new_mems[index];
                        mem->flags = new_class;
                        mem->score += 2;
                    }
                }
                if (cutoff < num_entries) {
                    std::lock_guard lock(memories_mutex);
                    for (size_t i = cutoff; i < num_entries; i++) {
                        auto [ip, mem] = updated_memories[i];
                        memories[ip] = mem;
                    }
                }
                break;
            }
        }
    }
}

extern "C" int initialize(bool activate_sniffer) {
    try {
        forgery_model = torch::jit::load(FORGERY_MODEL_PATH);
        packet_model = torch::jit::load(PACKET_MODEL_PATH);
        merge_model = torch::jit::load(MERGE_MODEL_PATH);
        packet_cls_model = torch::jit::load(PACKET_CLS_MODEL_PATH);
        device_cls_model = torch::jit::load(DEVICE_CLS_MODEL_PATH);
    }
    catch (const c10::Error& e) {
        return 1;
    }
    try {
        std::random_device true_rng;
        rng.seed(true_rng());
    } catch (...) {
        rng.seed(std::chrono::system_clock::now().time_since_epoch().count());
    }
    limnet_thread = std::thread(limnet_thread_func);
    return 0;
}

extern "C" void terminate() {
    needs_to_terminate = true;
    limnet_thread.join();
}

extern "C" void on_memories_received(ip_t* ips, void* memories, size_t num_entries) {
    torch::NoGradGuard no_grad;
    auto compressed_memories = torch::from_blob(memories, {(int64_t)num_entries, (int64_t)MEMORY_DIM}, c10::TensorOptions().dtype(COMPRESSED_DTYPE));
    auto decompressed_memories = compressed_memories.to(c10::kFloat);
    auto has_forgeries = forgery_model.forward({decompressed_memories}).toBool();
    if (!has_forgeries)
        message_queue.enqueue(message_t(std::in_place_type<shared_mem_t>, std::vector(ips, ips + num_entries), decompressed_memories));
}

extern "C" void on_packet_received(ip_t ip_src, ip_t ip_dst, float* features) {
    auto packet = std::make_unique<packet_t>(ip_src, ip_dst);
    memcpy(&packet->features, features, sizeof(packet_t::features));
    message_queue.enqueue(message_t(std::in_place_type<std::unique_ptr<packet_t>>, std::move(packet)));
}

extern "C" size_t get_memories_to_share(ip_t* ips, void* mems, size_t num_entries) {
    auto size = compressed_memory_size();
    std::vector<std::pair<ip_t, memory_t> > top_mems(num_entries);
    {
        std::lock_guard lock(memories_mutex);
        num_entries = std::distance(top_mems.begin(), std::partial_sort_copy(memories.begin(), memories.end(), top_mems.begin(), top_mems.end(), [](const auto& a, const auto& b){ return a.second->score < b.second->score; }));
    }
    for (size_t i = 0; i < num_entries; i++) {
        torch::Tensor compressed_memory;
        {
            std::lock_guard lock(top_mems[i].second->mutex);
            top_mems[i].second->score = 0;
            compressed_memory = top_mems[i].second->memory.to(c10::kHalf);
        }
        memcpy(mems, compressed_memory.data_ptr(), size);
        mems = (uint8_t*)mems + size;
        ips[i] = top_mems[i].first;
    }
    return num_entries;
}

extern "C" ip_t get_random_peer() {
    std::lock_guard lock(memories_mutex);
    auto mem_iter = memories.begin();
    std::advance(mem_iter, std::uniform_int_distribution((size_t)0, memories.size() - 1)(rng));
    return mem_iter->first;
}

extern "C" size_t compressed_memory_size() {
    return MEMORY_DIM * c10::elementSize(COMPRESSED_DTYPE);
}

extern "C" size_t packet_features_size() {
    return sizeof(packet_t::features);
}