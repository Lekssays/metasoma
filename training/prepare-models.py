import os
import sys
import torch


class ForgeryModel(torch.nn.Module):
    def __init__(self, flow):
        super().__init__()
        self.flow = flow
        
    def forward(self, recv_mems: torch.Tensor):
        return torch.max(0.5 - self._log_prob(recv_mems, None)).sign().clamp(0, 1).to(torch.bool)
        
    def _log_prob(self, inputs, context):
        embedded_context = self.flow._embedding_net(context)
        noise, logabsdet = self.flow._transform(inputs, context=embedded_context)
        log_prob = self.flow._distribution._log_prob(noise, context=embedded_context)
        return log_prob + logabsdet

class PacketModel(torch.nn.Module):
    def __init__(self, incoming: torch.nn.Module, outgoing: torch.nn.Module):
        super().__init__()
        self.incoming = incoming
        self.outgoing = outgoing
        
    def forward(self, src_mem: torch.Tensor, dst_mem: torch.Tensor, packet_feats: torch.Tensor):
        return (self.incoming(torch.cat((src_mem, packet_feats)).unsqueeze(0), dst_mem.unsqueeze(0)).squeeze(),
                self.outgoing(torch.cat((dst_mem, packet_feats)).unsqueeze(0), src_mem.unsqueeze(0)).squeeze())


class PacketClsModel(torch.nn.Module):
    def __init__(self, cls: torch.nn.Module):
        super().__init__()
        self.cls = cls
        
    def forward(self, src_mem: torch.Tensor, dst_mem: torch.Tensor, packet_feats: torch.Tensor):
        score = self.cls(torch.cat((src_mem, dst_mem, packet_feats)).unsqueeze(0)).squeeze()
        return score.sign().clamp(0, 1).to(torch.uint8)


class DeviceClsModel(torch.nn.Module):
    def __init__(self, cls: torch.nn.Module):
        super().__init__()
        self.cls = cls
        
    def forward(self, device_mem: torch.Tensor):
        scores = self.cls(device_mem.unsqueeze(0)).squeeze()
        scores = scores.sign().clamp(0, 1).to(torch.uint8)
        return scores[1] * 2 + scores[0]


class MergeModel(torch.nn.Module):
    def __init__(self, merger: torch.nn.Module):
        super().__init__()
        self.merger = merger
        
    def forward(self, old_mems: torch.Tensor, new_mems: torch.Tensor, indices: torch.Tensor):
        return self.merger(old_mems, new_mems[indices])


def test_models_before_scripting(trainer):
    B = 5
    mem = torch.rand(trainer.Dmem)
    feats = torch.rand(trainer.updater.incoming.input_size - trainer.Dmem)
    mems = torch.rand(B, trainer.Dmem)
    ids = torch.arange(B)

    packet_model = torch.jit.trace(PacketModel(trainer.updater.incoming, trainer.updater.outgoing), (mem, mem, feats))
    packet_cls_model = torch.jit.trace(PacketClsModel(trainer.pkt_cls), (mem, mem, feats))
    device_cls_model = torch.jit.trace(DeviceClsModel(trainer.dev_cls), mem)
    merge_model = torch.jit.trace(MergeModel(trainer.merger), (mems, mems, ids))
    forgery_model = torch.jit.trace(ForgeryModel(flow), mems)

    B = 6
    mems = torch.rand(B, trainer.Dmem)
    ids = torch.arange(B)

    src_out, dst_out = packet_model(mem, mem, feats)
    assert src_out.shape == (trainer.Dmem,)
    assert dst_out.shape == (trainer.Dmem,)
    assert packet_cls_model(mem, mem, feats).shape == torch.Size([])
    assert packet_cls_model(mem, mem, feats).dtype == torch.uint8
    assert device_cls_model(mem).shape == torch.Size([])
    assert device_cls_model(mem).dtype == torch.uint8
    assert merge_model(mems, mems, ids).shape == (B, trainer.Dmem)
    assert forgery_model(mems).shape == torch.Size([])
    assert forgery_model(mems).dtype == torch.bool

        
if __name__ == '__main__':
    root_folder = sys.argv[1]

    print('### Loading ###')
    trainer = torch.load(os.path.join(root_folder, 'trainer.pt'))
    flow = torch.load(os.path.join(root_folder, 'flow.pt'))

    print('### Testing ###')
    test_models_before_scripting(trainer)

    print('### Scripting ###')
    os.makedirs('deployment/generated', exist_ok = True)
    torch.jit.trace(ForgeryModel(flow), torch.randn(5, trainer.Dmem)).save('deployment/generated/forgery_model.pt')
    torch.jit.script(PacketModel(trainer.updater.incoming, trainer.updater.outgoing)).save('deployment/generated/packet_model.pt')
    torch.jit.script(PacketClsModel(trainer.pkt_cls)).save('deployment/generated/packet_cls_model.pt')
    torch.jit.script(DeviceClsModel(trainer.dev_cls)).save('deployment/generated/device_cls_model.pt')
    torch.jit.script(MergeModel(trainer.merger)).save('deployment/generated/merge_model.pt')
    
    print('### Saving info ###')
    os.makedirs('gossip/generated', exist_ok = True)
    with open('gossip/training_info.h.template', 'r') as template_file:
        template = template_file.read()
    with open('gossip/generated/training_info.h', 'w') as generated_header:
        generated_header.write(template.format(
            memory_dim = trainer.Dmem,
            num_protos = trainer.updater.incoming.input_size - trainer.Dmem - 5,
            forgery_model_path = 'forgery_model.pt',
            packet_model_path = 'packet_model.pt',
            packet_cls_model_path = 'packet_cls_model.pt',
            device_cls_model_path = 'device_cls_model.pt',
            merge_model_path = 'merge_model.pt'
        ))

    print('### Done ###')