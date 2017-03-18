FROM ubuntu:xenial

# RUN COMMAND

RUN mkdir -p /var/www/go/bin; \
    mkdir libgit2

RUN apt-get -y update; \
    apt-get -y upgrade; \
    apt-get -y install cmake; \
    apt-get -y install git; \
    apt-get clean all

# Add Compiled Gogeta 
COPY Gogeta /var/www/go/bin

ENV GOOGLE_APPLICATION_CREDENTIALS="/usr/local/gcloud-service-key.json"
ENV PAPERTRAIL_ENDPOINT="logs5.papertrailapp.com:54016"
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

# Add libgit2
COPY libgit2 /libgit2
RUN cd libgit2; \
    rm -rf build/; \
    mkdir build; \
    cd build; \
    cmake .. -DBUILD_CLAR=OFF -DCMAKE_INSTALL_PREFIX=/usr/local; \
    cmake --build . --target install; \
    rm -rf /libgit2

# Run Gogeta
CMD echo $GCLOUD_SERVICE_KEY | base64 --decode --ignore-garbage > /usr/local/gcloud-service-key.json; \
    /var/www/go/bin/Gogeta
