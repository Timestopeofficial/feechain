# Build
FROM golang:1.16 AS build

RUN apt update -y && \
    apt install libgmp-dev libssl-dev git -y && \
    apt -y clean all

ENV GOPATH=/root/go
ENV GO111MODULE=on
ENV HMY_PATH=${GOPATH}/src/github.com/Timestopeofficial
RUN mkdir -p $HMY_PATH

WORKDIR $HMY_PATH

RUN git clone https://github.com/Timestopeofficial/feechain.git && \
    git clone https://github.com/Timestopeofficial/bls.git && \
    git clone https://github.com/Timestopeofficial/mcl.git

WORKDIR $HMY_PATH/feechain

RUN go mod download

RUN go mod tidy

RUN make linux_static && \
    cp ./bin/feechain /root/feechain && \
    cp ./rosetta/infra/run.sh /root/run.sh && \
    cp ./rosetta/infra/rclone.conf /root/rclone.conf

RUN cp ./rosetta/infra/feechain-pstn.conf /root/feechain-pstn.conf && \
    cp ./rosetta/infra/feechain-mainnet.conf /root/feechain-mainnet.conf && \
    cp ./.hmy/rosetta_local_fix.csv /root/rosetta_local_fix.csv

# Execution
FROM ubuntu:20.04

RUN apt update -y && \
    apt install libgmp-dev libssl-dev ca-certificates rclone -y && \
    apt -y clean all \

WORKDIR /root

COPY --from=build /root/feechain /root/feechain
COPY --from=build /root/run.sh /root/run.sh
COPY --from=build /root/rclone.conf /root/.config/rclone/rclone.conf
COPY --from=build /root/feechain-pstn.conf /root/feechain-pstn.conf
COPY --from=build /root/feechain-mainnet.conf /root/feechain-mainnet.conf
COPY --from=build /root/rosetta_local_fix.csv /root/rosetta_local_fix.csv

EXPOSE 9500/tcp
EXPOSE 9700/tcp
EXPOSE 9800/tcp
ENTRYPOINT ["/bin/bash","/root/run.sh"]
