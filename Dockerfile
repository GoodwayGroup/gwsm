FROM alpine:3.14.0

COPY gwsm /usr/local/bin/gwsm
RUN chmod +x /usr/local/bin/gwsm

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/gwsm" ]