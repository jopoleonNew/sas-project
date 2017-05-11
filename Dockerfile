FROM alpine:latest                                                                                                                                                                
COPY ./static /app/static
COPY ./conf-docker.json /app/configuration.json
COPY ./sas /app/sas
RUN apk --no-cache add ca-certificates && update-ca-certificates \
        && chown -R root:root /app \
            && chmod +x /app/sas \
                && chmod -R 770 /app
EXPOSE 3000
WORKDIR /app
CMD ["/app/sas", "-config", "/app/configuration.json"]