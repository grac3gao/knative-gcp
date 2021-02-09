FROM ubuntu:16.04
RUN apt-get update \
    && apt-get install -y mysql-client \
        telnet \
        net-tools \
        curl \
        iputils-ping \
        netcat \
        redis-tools \
        default-jre
CMD ["/bin/bash", "-c", "tail -f /dev/null"]

