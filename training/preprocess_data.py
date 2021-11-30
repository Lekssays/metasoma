import glob
import ipaddress
import os
import pandas as pd
from tqdm.auto import tqdm

def extract_features(input_paths, output_path):
    columns = {'ip_src', 'ip_dst', 'timestamp', 'length', 'protocol'}

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

    data = data[['timestamp', 'ip_src_int', 'ip_dst_int', 'src_is_private', 'src_is_multicast', 'dst_is_private', 'dst_is_multicast', 'length'] + list(one_hot_protocols.columns)]

    os.makedirs(output_path, exist_ok = True)
    private_ips = map(lambda ip: (int(ip), str(ip)), filter(lambda ip: ip.is_private and not ip.is_multicast and int(ip) > 0 and int(ip) < 2**32 - 1, pd.unique(ips.values.ravel('K'))))
    for ip_int, ip_str in tqdm(list(private_ips), desc = 'writing data'):
        data[(data.ip_src_int == ip_int) | (data.ip_dst_int == ip_int)].to_csv(f'{output_path}/{ip_str}.csv', index = False, header = False)

    return data

if __name__ == '__main__':
    extract_features(['data/raw/kitsune/*'], 'data/processed/kitsune')
    # extract_features(['data/raw/medbiot/*'], 'data/processed/medbiot')
    # extract_features(['data/raw/medbiot/torii*'], 'data/processed/torii')
    # extract_features(['data/raw/medbiot/mirai*'], 'data/processed/mirai')
    # extract_features(['data/raw/medbiot/bashlite*'], 'data/processed/bashlite')
