FROM ubuntu:21.10

SHELL ["/bin/bash", "--login", "-c"]

RUN                                                                                                                         \
# 1) install Ubuntu packages
    apt update                                                                                                              \
 && DEBIAN_FRONTEND=noninteractive apt install -y libpcap-dev unzip build-essential wget git                                \
 && rm -rf /var/lib/apt/lists/*                                                                                             \
# 2) install PcapPlusPlus
 && wget -q https://github.com/seladb/PcapPlusPlus/releases/download/v21.05/pcapplusplus-21.05-ubuntu-20.04-gcc-9.tar.gz    \
 && tar -xf pcapplusplus-21.05-ubuntu-20.04-gcc-9.tar.gz                                                                    \
 && rm pcapplusplus-21.05-ubuntu-20.04-gcc-9.tar.gz                                                                         \
 && cd pcapplusplus-21.05-ubuntu-20.04-gcc-9                                                                                \
 && ./install.sh                                                                                                            \
 && cd /                                                                                                                    \
 && rm -rf pcapplusplus-21.05-ubuntu-20.04-gcc-9                                                                            \
# 3) install PyTorch C++
 && wget -q https://download.pytorch.org/libtorch/cpu/libtorch-cxx11-abi-shared-with-deps-1.9.1%2Bcpu.zip                   \
 && unzip libtorch-cxx11-abi-shared-with-deps-1.9.1+cpu.zip                                                                 \
 && rm libtorch-cxx11-abi-shared-with-deps-1.9.1+cpu.zip                                                                    \
 && cp libtorch/lib/* /usr/local/lib                                                                                        \
 && cp -r libtorch/include/* /usr/local/include                                                                             \
 && rm -rf libtorch                                                                                                         \
# 4) install ConcurrentQueue
 && wget -q https://github.com/cameron314/concurrentqueue/archive/refs/tags/v1.0.3.zip                                      \
 && unzip v1.0.3.zip                                                                                                        \
 && rm v1.0.3.zip                                                                                                           \
 && mkdir /usr/local/include/concurrentqueue                                                                                \
 && cp concurrentqueue-1.0.3/*.h concurrentqueue-1.0.3/*.md /usr/local/include/concurrentqueue                              \
 && rm -rf concurrentqueue-1.0.3                                                                                            \
# 5) install Go
 && wget -q https://golang.org/dl/go1.17.2.linux-amd64.tar.gz                                                               \
 && tar -C /usr/local -xzf go1.17.2.linux-amd64.tar.gz                                                                      \
# 6) install Miniconda
 && wget -q https://repo.anaconda.com/miniconda/Miniconda3-py39_4.10.3-Linux-x86_64.sh                                      \
 && chmod +x Miniconda3-py39_4.10.3-Linux-x86_64.sh                                                                         \
 && ./Miniconda3-py39_4.10.3-Linux-x86_64.sh -b -p /miniconda                                                               \
 && /miniconda/bin/conda init bash                                                                                          \
# 7) install Python packages
 && /miniconda/bin/conda install -q -y numpy pandas                                                                         \
 && /miniconda/bin/conda install -q -y pytorch cudatoolkit=11.3 -c pytorch                                                  \
 && /miniconda/bin/conda clean -a -q -y                                                                                     \
# 8) install redis
 && apt install redis-server systemctl                                                                                      \
 && sed -i 's/supervised no/supervised systemd/' /etc/redis/redis.conf                                                      \
 && systemctl restart redis.service

ENV PATH="${PATH}:/usr/local/go/bin"
ENV LD_LIBRARY_PATH="/usr/local/lib"