import glob
import ipaddress
import numpy as np
import os
import pandas as pd
import torch
from tqdm.auto import tqdm

class Window:
    def __init__(self, Nprivs, t, ids, Xe, Ye, Ymal, Yatk):
        self.Nprivs = Nprivs
        self.t = t
        self.ids = ids
        self.Xe = Xe
        self.Ye = Ye
        self.Ymal = Ymal
        self.Yatk = Yatk

    def __getitem__(self, idx):
        return Window(self.Nprivs, self.t[idx], self.ids[idx], self.Xe[idx], self.Ye[idx], self.Ymal[idx], self.Yatk[idx])

def pandas2torch(data, dtype = np.float32):
    return torch.from_numpy(dtype(data.to_numpy()))

def extract_features(input_paths, deployment_output_path, training_output_path):
    columns = {'ip_src', 'ip_dst', 'timestamp', 'length', 'protocol', 'traffic_type', 'port_src', 'port_dst'}

    files = tqdm([file for path in input_paths for file in glob.iglob(path)], desc = 'loading data')
    data = pd.concat((pd.read_csv(file, usecols = columns) for file in files), ignore_index = True, copy = False)

    data.sort_values('timestamp', kind = 'quicksort', inplace = True, ignore_index = True)

    ips = data[['ip_src', 'ip_dst']].applymap(ipaddress.IPv4Address)
    data['ip_src_int'] = ips['ip_src'].map(int)
    data['ip_dst_int'] = ips['ip_dst'].map(int)
    data['src_is_private'] = ips['ip_src'].map(lambda ip: int(ip.is_private))
    data['dst_is_private'] = ips['ip_dst'].map(lambda ip: int(ip.is_private))
    data['src_is_multicast'] = ips['ip_src'].map(lambda ip: int(ip.is_multicast))
    data['dst_is_multicast'] = ips['ip_dst'].map(lambda ip: int(ip.is_multicast))

    q1, med, q9 = data['length'].quantile([.1,.5,.9])
    data['length'] = (data['length'] - med) / (q9 - q1)
    one_hot_protocols = pd.get_dummies(data['protocol'])
    data = data.join(one_hot_protocols)

    features = ['src_is_private', 'src_is_multicast', 'dst_is_private', 'dst_is_multicast', 'length'] + list(one_hot_protocols.columns)
    data_out = data[['timestamp', 'ip_src_int', 'ip_dst_int'] + features]

    unique_ips = pd.unique(ips.values.ravel('K'))
    is_private = lambda ip: ip.is_private and not ip.is_multicast and int(ip) > 0 and int(ip) < 2**32 - 1
    private_ips = list(map(lambda ip: (int(ip), str(ip)), filter(is_private, unique_ips)))

    os.makedirs(deployment_output_path, exist_ok = True)
    for ip_int, ip_str in tqdm(private_ips, desc = 'writing data'):
        data_out[(data_out.ip_src_int == ip_int) | (data_out.ip_dst_int == ip_int)].to_csv(f'{deployment_output_path}/{ip_str}.csv', index = False, header = False)

    sorted_ips = [int(ip) for ip in unique_ips if is_private(ip)]
    Nprivs = len(sorted_ips)
    sorted_ips = sorted_ips + [int(ip) for ip in unique_ips if not is_private(ip)]
    ip2id = {ip: id for id, ip in enumerate(sorted_ips)}

    t = pandas2torch(data.timestamp)
    ids = torch.empty((t.shape[0], 2), dtype = torch.int64)
    ids[:,0] = pandas2torch(data.ip_src_int.map(lambda ip: ip2id[ip]), dtype = np.int64)
    ids[:,1] = pandas2torch(data.ip_dst_int.map(lambda ip: ip2id[ip]), dtype = np.int64)
    Xe = pandas2torch(data[features])
    Ye = data['traffic_type'] == 'mal'
    Ymal = torch.empty(2 * t.shape[0])
    Ymal[0::2] = pandas2torch(Ye & (data['port_src'] > 10000))
    Ymal[1::2] = pandas2torch(Ye & (data['port_dst'] > 10000))
    Yatk = torch.empty(2 * t.shape[0])
    Yatk[0::2] = pandas2torch(Ye & (data['port_src'] < 10000))
    Yatk[1::2] = pandas2torch(Ye & (data['port_dst'] < 10000))
    os.makedirs(os.path.dirname(training_output_path), exist_ok = True)
    torch.save(Window(Nprivs, t, ids, Xe, pandas2torch(Ye), Ymal, Yatk), training_output_path)

    return data

if __name__ == '__main__':
    extract_features(['data/raw/kitsune/*'], 'data/deployment/kitsune', 'data/training/kitsune.pt')
    # extract_features(['data/raw/medbiot/*'], 'data/processed/medbiot')
    # extract_features(['data/raw/medbiot/torii*'], 'data/processed/torii')
    # extract_features(['data/raw/medbiot/mirai*'], 'data/processed/mirai')
    # extract_features(['data/raw/medbiot/bashlite*'], 'data/processed/bashlite')