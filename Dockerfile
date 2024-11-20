FROM alpine:latest
LABEL authors="wangao"
USER root
WORKDIR /JUGGERNAUT
COPY juggernaut config.yaml ./
RUN  chmod +x juggernaut && mkdir config && mv config.yaml /JUGGERNAUT/config/
EXPOSE 8080
ENTRYPOINT ["./juggernaut"]