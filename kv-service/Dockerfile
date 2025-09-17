FROM --platform=linux/amd64 debian:bookworm-slim
LABEL "maintainer"="Brendan Wilson <bmwilson.dev@gmail.com"
LABEL "version"="0.1"
LABEL "description"="Censys KV Store Server"
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
# copy in application
RUN git clone git@github.com:lossdev/censys-kv.git /opt/censys-kv
# build
WORKDIR /opt/censys-kv/kv-service
RUN /usr/local/go/bin/go mod tidy 
RUN /usr/local/go/bin/go build -o ./kv-server ./main.go
RUN chown -R kv /opt/censys-kv
# run as kv
ENTRYPOINT ["sudo", "-u", "kv", "/opt/censys-kv/kv-service/kv-server"]
