FROM envoyproxy/envoy:v1.11.2
RUN apt-get -qq update
COPY config.yaml /etc/envoy.yaml
CMD /usr/local/bin/envoy -c /etc/envoy.yaml
