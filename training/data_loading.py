import heapq
import random
from re import I
import torch
from tqdm import trange
from models import PACKET, GOSSIP

# https://docs.python.org/3/library/heapq.html#priority-queue-implementation-notes
from dataclasses import dataclass, field
from typing import Any
@dataclass(order = True)
class Trigger:
    id: int = field(compare = False)
    t: int

def uniform(start, stop):
    return random.uniform(start, stop)

def after(start, delay):
    return start + random.uniform(delay * 0.85, delay * 1.15)

def make_batch(data, pos, avg_gossip_secs, batch_size, num_send):
    batch = []
    bYe = []
    bYmal = []
    bYatk = []

    tstart = data.t[pos].item()
    triggers = [Trigger(i, uniform(tstart, tstart + avg_gossip_secs)) for i in range(data.Nprivs)]
    heapq.heapify(triggers)

    seen = [set([i]) for i in range(data.Nprivs)]
    devs = list(range(data.Nprivs))

    for pos in range(pos, pos + batch_size):

        while triggers[0].t < data.t[pos].item():
            srcdev = triggers[0].id
            dstdevs = list(dev for dev in seen[srcdev] if dev < data.Nprivs and dev != srcdev)
            if len(dstdevs) == 0:
                dstdevs = devs

            dstdev = srcdev
            while dstdev == srcdev:
                dstdev = random.choice(dstdevs)

            to_send = torch.tensor(random.sample(devs, num_send))
            batch.append((GOSSIP, (srcdev, dstdev, to_send)))
            triggers[0].t = after(triggers[0].t, avg_gossip_secs)
            heapq.heapreplace(triggers, triggers[0])

        privs = torch.tensor([id.item() for id in data.ids[pos] if id < data.Nprivs])
        if len(privs) > 0:

            src, dst = data.ids[pos]
            if src < data.Nprivs:
                seen[src].add(dst)
            if dst < data.Nprivs:
                seen[dst].add(src)

            batch.append((PACKET, (privs, src, dst, data.Xe[pos])))

            for _ in range(len(privs)):
                bYe.append(data.Ye[pos])
                bYmal.append(data.Ymal[2 * pos])
                bYmal.append(data.Ymal[2 * pos + 1])
                bYatk.append(data.Yatk[2 * pos])
                bYatk.append(data.Yatk[2 * pos + 1])

    bYe = torch.tensor(bYe)
    bYmal = torch.tensor(bYmal)
    bYatk = torch.tensor(bYatk)
    return batch, (bYe, bYmal, bYatk)

class Dataset(torch.utils.data.Dataset):
    def __init__(self, data, windows, avg_gossip_secs, num_send, batch_size):
        super().__init__()
        self.data = data
        self.windows = windows
        self.avg_gossip_secs = avg_gossip_secs
        self.num_send = num_send
        self.batch_size = batch_size

    def __len__(self):
        return len(self.windows)

    def __getitem__(self, idx):
        pos = self.windows[idx]
        return make_batch(self.data, pos, self.avg_gossip_secs,
                          self.batch_size, self.num_send)
