ARG ARCH=aarch64
ARG VERSION=1.13
ARG UBUNTU_VERSION=22.04
ARG REPO=axisecp
ARG SDK=acap-native-sdk

FROM ${REPO}/${SDK}:${VERSION}-${ARCH}-ubuntu${UBUNTU_VERSION}

ARG ARCH
RUN echo "Architecture is: ${ARCH}"
ARG SDK_LIB_PATH_BASE=/opt/axis/acapsdk/sysroots/${ARCH}/usr
RUN echo "SDK library path base is: ${SDK_LIB_PATH_BASE}"

# Install Meson, Ninja, and other build dependencies
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ninja-build \
    git \
    curl

RUN apt-get clean && rm -rf /var/lib/apt/lists/*
RUN pip3 install meson>=1.1

#-------------------------------------------------------------------------------
# Golang build
#-------------------------------------------------------------------------------
ARG GOLANG_VERSION=1.22.0
RUN curl -fsSL "https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" -o golang.tar.gz \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:/usr/local/go/bin:${PATH}"
RUN mkdir -p "${GOPATH}/src" "${GOPATH}/bin" "${GOPATH}/pkg" \
    && chmod -R 777 "${GOPATH}"
### Golang End

#-------------------------------------------------------------------------------
# Gstreamer build
#-------------------------------------------------------------------------------
ARG BUILD_DIR=/opt/build
ARG GSTREAMER_VER=1.18.5
WORKDIR ${BUILD_DIR}
RUN git clone https://gitlab.freedesktop.org/gstreamer/gstreamer.git
WORKDIR ${BUILD_DIR}/gstreamer
ARG CROSS_FILE=
COPY ./meson/${CROSS_FILE} ${BUILD_DIR}/gstreamer
RUN apt-get update && apt-get install -y gcc g++ flex bison libfontconfig1-dev libdrm-dev --no-install-recommends 
ARG GST_NOBUILD_MAIN="-Drtsp_server=disabled -Dlibav=disabled -Ddevtools=disabled -Dgst-examples=disabled -Dpython=disabled -Dtests=disabled -Dtools=disabled -Dexamples=disabled -Ddoc=disabled"
ARG GST_NOBUILD_BASE_PLUGS="-Dgst-plugins-base:alsa=disabled -Dgst-plugins-base:ogg=disabled -Dgst-plugins-base:vorbis=disabled -Dgst-plugins-base:opus=disabled -Dgst-plugins-base:pango=disabled"
ARG GST_NOBUILD_GOOD_PLUGS="-Dgst-plugins-good:vpx=disabled -Dgst-plugins-good:isomp4=disabled -Dgst-plugins-good:cairo=disabled -Dgst-plugins-good:dv=disabled -Dgst-plugins-good:flac=disabled -Dgst-plugins-good:lame=disabled -Dgst-plugins-good:png=disabled"
ARG GST_NOBUILD_BAD_PLUGS="-Dgst-plugins-bad:nvcodec=disabled -Dgst-plugins-bad:qsv=disabled -Dgst-plugins-bad:aja=disabled -Dgst-plugins-bad:openjpeg=disabled -Dgst-plugins-bad:fdkaac=disabled -Dgst-plugins-bad:microdns=disabled -Dgst-plugins-bad:avtp=disabled -Dgst-plugins-bad:openh264=disabled"
# Setup the Meson build directory
RUN meson setup builddir \
        --cross-file=${CROSS_FILE} \
        --prefix=${SDK_LIB_PATH_BASE} \
        --libdir=${SDK_LIB_PATH_BASE}/lib \
        ${GST_NOBUILD_MAIN} \
        ${GST_NOBUILD_BASE_PLUGS} \
        ${GST_NOBUILD_GOOD_PLUGS} \
        ${GST_NOBUILD_BAD_PLUGS}

# Build and install using Ninja
RUN ninja -C builddir && \
    ninja -C builddir install

#-------------------------------------------------------------------------------
# Perpare the ACAP Build
#-------------------------------------------------------------------------------
ARG APP_NAME=app
ARG GO_ARCH=arm64
ARG GO_ARM=
ARG IP_ADDR= 
ARG PASSWORD= 
ARG START=
ARG INSTALL=
ARG GO_APP=test
ENV GO_APP=${GO_APP}
COPY . /opt/goaxis/
WORKDIR /opt/goaxis/${GO_APP}

#-------------------------------------------------------------------------------
# Copy gstreamer into app directory
#-------------------------------------------------------------------------------
RUN mkdir -p lib && \
    cp -a ${SDK_LIB_PATH_BASE}/lib/. lib/

#-------------------------------------------------------------------------------
# Perpare final build of golang acap app, check for install and start 
#-------------------------------------------------------------------------------
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=${GO_ARCH}
ENV GOARM=${GO_ARM}
ENV APP_NAME=${APP_NAME}
RUN . /opt/axis/acapsdk/environment-setup* && \
    unset PKG_CONFIG_SYSROOT_DIR && \
    printenv && \
    go build -ldflags "-s -w  -extldflags '-L./lib -Wl,-rpath,./lib'" -o ${APP_NAME} . && \
    acap-build --build no-build ./ && \
    if [ "${INSTALL}" = "YES" ]; then eap-install.sh ${IP_ADDR} ${PASSWORD} install; fi && \
    if [ "${START}" = "YES" ]; then eap-install.sh start; fi

#-------------------------------------------------------------------------------
# Create output directory, we copy files from eap to host
#-------------------------------------------------------------------------------
RUN mkdir /opt/eap
RUN mv *.eap /opt/eap
RUN ls /opt/eap

