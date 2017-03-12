FROM debian:squeeze
FROM golang:1.8.0

# Get cmake
RUN mkdir /libgit2; \
    apt-get -y update; \
    apt-get -y upgrade; \
    apt-get -y install cmake; \
    apt-get clean all

# Copy libgit2 file
COPY libgit2 /libgit2

# Install libgit2
RUN mkdir /libgit2/build; \
    cd /libgit2/build; \
    cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local; \
    cmake --build . --target install; \
    rm -R /libgit2
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
ENV PATH=$PATH:/go/bin

# Install Gogeta
RUN go get github.com/Gamebuildr/Gogeta

# Run Gogeta
CMD cd /go/src/github.com/Gamebuildr/Gogeta; \
    Gogeta

