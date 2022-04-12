import os
from posixpath import split
import torch
from tqdm import tqdm
from torch.nn.functional import binary_cross_entropy_with_logits as bcelogits

from models import Trainer
from data_loading import Dataset
from data_preparation import Window

def loss_fn(outs, labels):
    Ze, Zn = outs
    Ye, Ymal, Yatk = labels
    return bcelogits(Ze, Ye) + bcelogits(Zn[:,0], Ymal) + bcelogits(Zn[:,1], Yatk)

def print_accuracy(outs, labels, epoch):
    Ze, Zn = outs
    Ye, Ymal, Yatk = labels
    acc = lambda z, y: torch.sum(torch.round(torch.sigmoid(z)) == y).item() / len(z)
    print(f'### Epoch {epoch} ###')
    print(f'edge-mal-acc: {round(acc(Ze, Ye), 3)}')
    print(f'node-mal-acc: {round(acc(Zn[:,0], Ymal), 3)}')
    print(f'node-atk-acc: {round(acc(Zn[:,1], Yatk), 3)}')
    print('')

def split_windows(tot, size, shift, p):
    windows = torch.randperm((tot - size) // shift) * shift
    cutoff = int(len(windows) * (1 - p))
    return windows[:cutoff], windows[cutoff:]

if __name__ == '__main__':
    data = torch.load('data/training/kitsune.pt')
    train_windows, valid_windows = split_windows(data.t.shape[0], 5000, 500, 0.15)

    train_dataset = Dataset(data, train_windows, 15, 8, 5000)
    train_dataloader = torch.utils.data.DataLoader(train_dataset, shuffle = True, batch_size = None, num_workers = 2, persistent_workers = True)

    device = 'cpu'
    root_folder = 'runs/test-run'
    gi = 0
    os.makedirs(root_folder, exist_ok = True)

    model = Trainer(train_dataset.data.Xe.shape[1], 32, 2, 1).to(device)
    optimizer = torch.optim.Adam(model.parameters())

    for epoch in range(5):

        model.train()
        for batch, labels in (progr := tqdm(train_dataloader, desc = 'train')):
            H = torch.zeros((train_dataset.data.Nprivs, torch.max(train_dataset.data.ids) + 1, model.Dmem), device = device)
            (gsp_srcs, gsp_tgts), outs = model(H, batch)
            loss = loss_fn(outs, labels)
            loss.mean().backward()
            optimizer.step()
            optimizer.zero_grad()
            progr.set_description(f'train (loss: {round(loss.item(), 3)})')
            torch.save((labels[1], gsp_srcs, gsp_tgts), f'{root_folder}/gossip-{gi}.pt')
            gi += 1

        torch.save(model, f'{root_folder}/trainer.pt')

        model.eval()
        with torch.no_grad():
            Ze = []
            Zn = []
            Ye = []
            Ymal = []
            Yatk = []
            for batch, (ye, ymal, yatk) in tqdm(train_dataloader, desc = 'valid'):
                H = torch.zeros((train_dataset.data.Nprivs, torch.max(train_dataset.data.ids) + 1, model.Dmem))
                ze, zn = model(H, batch)
                Ze.append(ze)
                Zn.append(zn)
                Ye.append(ye)
                Ymal.append(ymal)
                Yatk.append(yatk)
            Ze = torch.cat(Ze)
            Zn = torch.cat(Zn)
            Ye = torch.cat(Ye)
            Ymal = torch.cat(Ymal) 
            Yatk = torch.cat(Yatk)
            print_accuracy((Ze, Zn), (Ye, Ymal, Yatk), epoch)

