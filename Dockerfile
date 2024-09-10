FROM scratch
COPY bin/sml /sml
ENTRYPOINT ["/sml"]
