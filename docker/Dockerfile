FROM alpine:3.1
ENV UNITY3D2PNG_PROXY_VERSION=v0.0.1
ENV WORKDIR=
ENV REPOSITORY=
ENV UNITY3D2PNG_URL=
ENV PATH_TEMPLATE=
RUN apk add --update curl bash git openssh-client
WORKDIR /root
RUN mkdir -p /root/.ssh
RUN touch /root/.ssh/id_rsa
RUN echo -e "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN curl -L -O https://github.com/mix3/unity3dpng-proxy/releases/download/${UNITY3D2PNG_PROXY_VERSION}/unity3d2png-proxy-${UNITY3D2PNG_PROXY_VERSION}-linux-amd64.zip && \
    unzip unity3d2png-proxy-${UNITY3D2PNG_PROXY_VERSION}-linux-amd64.zip unity3d2png-proxy/unity3d2png-proxy && \
    mv unity3d2png-proxy unity3d2png-proxy_dir && \
    mv unity3d2png-proxy_dir/unity3d2png-proxy unity3d2png-proxy && \
    rm -rf unity3d2png-proxy_dir/
CMD ["./unity3d2png-proxy"]
