FROM scratch
COPY stbuild/cloudtail /bin/cloudtail
ENTRYPOINT /bin/cloudtail
