FROM lekssays/declimnet

SHELL ["/bin/bash", "--login", "-c"]

RUN                                                                                                                         \
    apt update                                                                                                              \
    && wget -q https://repo.anaconda.com/miniconda/Miniconda3-py39_4.10.3-Linux-x86_64.sh                                      \
    && chmod +x Miniconda3-py39_4.10.3-Linux-x86_64.sh                                                                         \
    && ./Miniconda3-py39_4.10.3-Linux-x86_64.sh -b -p /miniconda                                                               \
    && /miniconda/bin/conda init bash                                                                                          \
    && /miniconda/bin/conda install -q -y numpy pandas                                                                         \
    && /miniconda/bin/conda install -q -y pytorch cudatoolkit=11.3 -c pytorch                                                  \
    && /miniconda/bin/conda clean -a -q -y                                                                                     \
    && systemctl enable redis-server                                                                                            \
    && service redis-server restart                                                                                              

ENV PATH="${PATH}:/usr/local/go/bin"
ENV LD_LIBRARY_PATH="/usr/local/lib"