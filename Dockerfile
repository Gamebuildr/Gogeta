FROM debian:squeeze

# Add Compiled Gogeta 
RUN mkdir -p /var/www/go/bin; \
    mkdir -p /var/www/go/big/client/logs
COPY Gogeta /var/www/go/bin

# Run Gogeta
CMD cd /var/www/go/bin; \
    Gogeta
