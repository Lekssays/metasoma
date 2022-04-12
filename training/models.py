import os
import torch
from math import prod
from tqdm.auto import tqdm

class MemoryUpdater(torch.nn.Module):
    def __init__(self, Din, Dout, cell_type):
        super().__init__()
        self.incoming = cell_type(Din + Dout, Dout)
        self.outgoing = cell_type(Din + Dout, Dout)

    def forward(self, Xe: torch.Tensor, Hsrc: torch.Tensor, Hdst: torch.Tensor):
        B = Hsrc.shape[:-1]
        srcInputs = torch.cat((Hdst, Xe), dim = -1).view(prod(B), -1)
        dstInputs = torch.cat((Hsrc, Xe), dim = -1).view(prod(B), -1)
        Hsrc = self.outgoing(srcInputs, Hsrc)
        Hdst = self.incoming(dstInputs, Hdst)
        return Hsrc.view(*B, -1), Hdst.view(*B, -1)

class MemoryMerger(torch.nn.Module):
    def __init__(self, D):
        super().__init__()
        self.linear = torch.nn.Linear(2 * D, D)

    def forward(self, Ha: torch.Tensor, Hb: torch.Tensor):
        return self.linear(torch.hstack((Ha, Hb)))

class MemoryMerger2(torch.nn.Module):
    def __init__(self, D, cell_type):
        super().__init__()
        self.cell = cell_type(D, D)

    def forward(self, Hsrc: torch.Tensor, Hdst: torch.Tensor):
        return self.cell(Hsrc, Hdst)

PACKET = 0
GOSSIP = 1 

class Trainer(torch.nn.Module):
    def __init__(self, Dfeat, Dmem, Ddev, Dpkt, *, cell_type = torch.nn.GRUCell):
        super().__init__()
        self.Dmem = Dmem
        self.updater = MemoryUpdater(Dfeat, Dmem, cell_type)
        self.merger = MemoryMerger2(Dmem, torch.nn.GRUCell)
        self.dev_cls = torch.nn.Linear(Dmem, Ddev)
        self.pkt_cls = torch.nn.Linear(2 * Dmem + Dfeat, Dpkt)

    def forward(self, H, steps):
        gsp_srcs = []
        gsp_tgts = []

        dev_outs = []
        pkt_outs = []
        for step in tqdm(steps, desc = 'sequence', leave = False):

            if step[0] == PACKET:
                Idev, Isrc, Idst, Xe = step[1]
                Xe = Xe.repeat(len(Idev), 1)
                Hsrc, Hdst = self.updater(Xe, H[Idev, Isrc], H[Idev, Idst])
                pkt_outs.append(self.pkt_cls(torch.cat((Hsrc, Hdst, Xe), dim = -1)))
                dev_outs.append(self.dev_cls(Hsrc))
                dev_outs.append(self.dev_cls(Hdst))
                H[Idev, Isrc] = Hsrc
                H[Idev, Idst] = Hdst

            elif step[0] == GOSSIP:
                Idevsrc, Idevdst, Isend = step[1]
                Hsrc = H[Idevsrc, Isend]
                Hdst = H[Idevdst, Isend]
                gsp_srcs.append(Hsrc.detach())
                gsp_tgts.append(Hdst.detach())
                H[Idevdst, Isend] = self.merger(Hsrc, Hdst)

        gsp_srcs = torch.cat(gsp_srcs)
        gsp_tgts = torch.cat(gsp_tgts)
        pkt_outs = torch.cat(pkt_outs).squeeze()
        dev_outs = torch.cat(dev_outs).squeeze()

        return (gsp_srcs, gsp_tgts), (pkt_outs, dev_outs)