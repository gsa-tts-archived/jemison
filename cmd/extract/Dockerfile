#FROM golang:bookworm
#FROM cloudfoundry/cflinuxfs4
FROM jemison/builder

# Copy in the binaries
WORKDIR /home/vcap/app
COPY . .
WORKDIR /home/vcap/app/cmd/extract
RUN make container_build
RUN chmod 755 /home/vcap/app/cmd/extract/run.sh

CMD ["/home/vcap/app/cmd/extract/run-fast.sh"]


