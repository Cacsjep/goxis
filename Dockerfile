ARG ARCH=aarch64
ARG VERSION=1.13
ARG UBUNTU_VERSION=22.04
ARG REPO=axisecp
ARG SDK=acap-native-sdk

FROM ${REPO}/${SDK}:${VERSION}-${ARCH}-ubuntu${UBUNTU_VERSION}

RUN apt-get update && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/*

### Golang
ARG GOLANG_VERSION=1.22.0
RUN curl -fsSL "https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" -o golang.tar.gz \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:/usr/local/go/bin:${PATH}"
RUN mkdir -p "${GOPATH}/src" "${GOPATH}/bin" "${GOPATH}/pkg" \
    && chmod -R 777 "${GOPATH}"
### Golang End

# Building the ACAP application
ARG APP_NAME=app
ARG GO_ARCH=arm64
ARG GO_ARM=
ARG IP_ADDR= 
ARG PASSWORD= 
ARG START=
ARG INSTALL=
COPY . /opt/goaxis/
WORKDIR /opt/goaxis/app
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=${GO_ARCH}
ENV GOARM=${GO_ARM}
ENV APP_NAME=${APP_NAME}
RUN echo ${APP_NAME}
RUN echo ${PASSWORD}
RUN . /opt/axis/acapsdk/environment-setup* && \
    go build -ldflags "-s -w" -o ${APP_NAME} . && \
    acap-build --build no-build ./ && \
    if [ "${INSTALL}" = "YES" ]; then eap-install.sh ${IP_ADDR} ${PASSWORD} install; fi && \
    if [ "${START}" = "YES" ]; then eap-install.sh start; fi
RUN mkdir /opt/eap
#RUN cp -r /opt/axis/acapsdk/sysroots/aarch64/usr/include/* /opt/eap/
RUN mv *.eap /opt/eap
RUN ls /opt/eap

