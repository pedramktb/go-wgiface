FROM golang:latest

RUN apt update -y && apt-get install -y iproute2

COPY . /wgiface
WORKDIR /wgiface

# Avoid loop on container test
RUN rm container_test.go

# Enable container only tests
RUN mv iface_tests.go iface_test.go

CMD ["sleep", "1000"]
