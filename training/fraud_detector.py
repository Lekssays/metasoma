import glob
import os
import sys
import torch

from tqdm.auto import tqdm, trange
from nflows.flows import MaskedAutoregressiveFlow

def load_memories(folder):
    batches = []
    for file in glob.iglob(f'{folder}/gossip-*.pt'):
        batches.append(torch.load(file)[1])
    return torch.cat(batches)

def fit_flow(data, flow, batch_size, epochs):
    dataset = torch.utils.data.TensorDataset(data)
    dataloader = torch.utils.data.DataLoader(dataset, batch_size, drop_last = True)
    optimizer = torch.optim.Adam(flow.parameters())

    for epoch in trange(1, epochs + 1, desc = 'Epochs'):
        for batch in (progress := tqdm(dataloader, total = len(data) // batch_size, desc = f'Epoch {epoch}')):
            loss = -flow.log_prob(batch[0]).mean()
            loss.backward()
            optimizer.step()
            optimizer.zero_grad()
            progress.set_description(f'Epoch {epoch} (loss: {loss.item()})')

if __name__ == '__main__':
    folder = sys.argv[1]
    memories = load_memories(folder)
    D = memories.shape[1]
    flow = MaskedAutoregressiveFlow(D, D, 2, 4)
    fit_flow(memories, flow, 64, 5)
    torch.save(flow, f'{folder}/flow.pt')

