FROM debian:latest
COPY dist/godynamicdns_0.0.1-next_linux_amd64.deb /godynamicdns_0.0.1-next_linux_amd64.deb
RUN apt update
RUN apt install /godynamicdns_0.0.1-next_linux_amd64.deb -y

FROM centos:latest
COPY dist/godynamicdns_0.0.1-next_linux_amd64.rpm /godynamicdns_0.0.1-next_linux_amd64.rpm
#RUN yum update
RUN dnf install -y /godynamicdns_0.0.1-next_linux_amd64.rpm