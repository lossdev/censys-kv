FROM --platform=linux/amd64 debian:bookworm-slim
LABEL "maintainer"="Brendan Wilson <bmwilson.dev@gmail.com"
LABEL "version"="0.1"
LABEL "description"="Censys KV Store Client"
# update, upgrade, set up go
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install sudo wget git openssh-client bash -y
RUN wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz -O /opt/go1.24.2.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf /opt/go1.24.2.linux-amd64.tar.gz
RUN rm /opt/go1.24.2.linux-amd64.tar.gz
# set up env variables, ssh and git for go mod tidy, and add kv user
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/bash.bashrc
RUN mkdir -p -m 0600 /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN useradd -r -s /bin/false kv
# copy in application, also have to cachebust this layer otherwise Docker gets confused
RUN git clone https://github.com/lossdev/censys-kv.git /opt/censys-kv2
# build
WORKDIR /opt/censys-kv2/kv-test-client
RUN /usr/local/go/bin/go mod tidy 
RUN /usr/local/go/bin/go build -o ./kv-client ./main.go
RUN chown -R kv /opt/censys-kv2
# run as kv
ENTRYPOINT ["sudo", "-u", "kv", "/opt/censys-kv2/kv-test-client/kv-client"]
