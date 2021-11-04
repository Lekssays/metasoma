import os
import torch

class ForgeryModel(torch.nn.Module):
    def __init__(self):
        super().__init__()
        
    def forward(recv_mems):
        return False

class PacketModel(torch.nn.Module):
    def __init__(self):
        super().__init__()
        
    def forward(src_mem, dst_mem, packet_feats):
        return src_mem, dst_mem

class PacketClsModel(torch.nn.Module):
    def __init__(self):
        super().__init__()
        
    def forward(src_mem, dst_mem, packet_feats):
        return 0

class DeviceClsModel(torch.nn.Module):
    def __init__(self):
        super().__init__()
        
    def forward(device_mem):
        return 0

class MergeModel(torch.nn.Module):
    def __init__(self):
        super().__init__()
        
    def forward(old_mems, new_mems, indices):
        return new_mems[indices]
        
if __name__ == '__main__':
    os.makedirs('deployment/generated', exist_ok = True)
    torch.jit.script(ForgeryModel()).save('deployment/generated/forgery_model.pt')
    torch.jit.script(PacketModel()).save('deployment/generated/packet_model.pt')
    torch.jit.script(PacketClsModel()).save('deployment/generated/packet_cls_model.pt')
    torch.jit.script(DeviceClsModel()).save('deployment/generated/device_cls_model.pt')
    torch.jit.script(MergeModel()).save('deployment/generated/merge_model.pt')
    
    os.makedirs('inference/generated', exist_ok = True)
    with open('inference/training_info.h.template', 'r') as template_file:
        template = template_file.read()
    with open('inference/generated/training_info.h', 'w') as generated_header:
        generated_header.write(template.format(
            memory_dim = 32,
            num_protos = 19,
            forgery_model_path = 'forgery_model.pt',
            packet_model_path = 'packet_model.pt',
            packet_cls_model_path = 'packet_cls_model.pt',
            device_cls_model_path = 'device_cls_model.pt',
            merge_model_path = 'merge_model.pt'
        ))