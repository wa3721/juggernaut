FROM alpine:latest
LABEL authors="wangao"
USER root
WORKDIR /JUGGERNAUT

# 复制juggernaut可执行文件和readme.md到工作目录
COPY juggernaut README.md ./

RUN mkdir -p config html && \
    chmod +x juggernaut \

COPY config.yaml config/

# 假设你的html文件在当前目录下的html文件夹中，将它们复制到镜像中的html目录
COPY html/ html/

# 暴露8080端口
EXPOSE 8080

# 设置容器启动时执行的命令
ENTRYPOINT ["./juggernaut"]