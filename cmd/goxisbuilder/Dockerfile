ARG ARCH=aarch64
ARG VERSION=1.13
ARG UBUNTU_VERSION=22.04
ARG REPO=axisecp
ARG SDK=acap-native-sdk 
FROM ${REPO}/${SDK}:${VERSION}-${ARCH}-ubuntu${UBUNTU_VERSION}

ARG ARCH
ARG VERSION
RUN echo "Architecture is: ${ARCH}"
ARG SDK_LIB_PATH_BASE=/opt/axis/acapsdk/sysroots/${ARCH}/usr
RUN echo "SDK library path base is: ${SDK_LIB_PATH_BASE}"
ARG APP_DIR=/opt/goaxis/
RUN mkdir ${APP_DIR}

# Install Meson, Ninja, and other build dependencies
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    wget \
    yasm \
    nasm \
    git \
    build-essential 

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

#-------------------------------------------------------------------------------
# Golang build
#-------------------------------------------------------------------------------
ARG GOLANG_VERSION=1.22.1
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
ARG FF_BUILD_DIR=/opt/build
ARG COMP_LIBAV=
ARG FFMPEG_VERSION=n5.1.2
ARG CROSS_PREFIX=
RUN mkdir -p ${FF_BUILD_DIR} && \
    if [ "${COMP_LIBAV}" = "YES" ]; then \
      git clone https://github.com/FFmpeg/FFmpeg.git ${FF_BUILD_DIR}/ffmpeg && \
      cd ${FF_BUILD_DIR}/ffmpeg && \
      git checkout ${FFMPEG_VERSION} && \
      ./configure \
        --arch=${ARCH} \
        --target-os=linux \
        --cross-prefix=${CROSS_PREFIX} \
        --enable-cross-compile \
        --prefix=${FF_BUILD_DIR} \
        --disable-everything \
        --disable-programs \
        --enable-avfilter \
        --enable-avformat \
        --enable-avcodec \
        --enable-avutil \
        --enable-parser=h264,mjpeg,hevc,aac,mp3 \
        --enable-bsf=h264_metadata,h264_mp4toannexb,hevc_metadata,hevc_mp4toannexb,aac_adtstoasc \
        --enable-protocol=file,rtmp,rtmpt,rtp,data,tcp,pipe,hls \
        --enable-encoder=h264,mjpeg,hevc \
        --enable-decoder=h264,mjpeg,hevc \
        --enable-muxer=flv,h264,mjpeg,hevc,mov,mpegts \
        --enable-gpl \
        --enable-small \
        --disable-doc \
        --enable-shared && \
      make && make install && \
      mkdir -p ${APP_DIR}/lib && \
      cp -a ${FF_BUILD_DIR}/lib/. ${APP_DIR}/lib; \
    fi


#-------------------------------------------------------------------------------
# Perpare the ACAP Build
#-------------------------------------------------------------------------------
ARG APP_NAME=app
ARG APP_MANIFEST=
ARG GO_ARCH=arm64
ARG GO_ARM=
ARG IP_ADDR= 
ARG PASSWORD= 
ARG START=
ARG INSTALL=
ARG GO_APP=test
ENV GO_APP=${GO_APP}
COPY . ${APP_DIR}
WORKDIR ${APP_DIR}
RUN python generate_makefile.py ${APP_NAME} ${GO_APP} ${APP_MANIFEST}
WORKDIR ${APP_DIR}/${GO_APP}

#-------------------------------------------------------------------------------
# Perpare final build of golang acap app, check for install and start 
#-------------------------------------------------------------------------------
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=${GO_ARCH}
ENV GOARM=${GO_ARM}
ENV APP_NAME=${APP_NAME}
ENV MANIFEST=${APP_MANIFEST}
RUN . /opt/axis/acapsdk/environment-setup* && \
    if [ "${COMP_LIBAV}" = "YES" ]; then \
        export CGO_LDFLAGS="-L${FF_BUILD_DIR}/lib/" && \
        export CGO_CFLAGS="-I${FF_BUILD_DIR}/include/" && \
        export PKG_CONFIG_PATH="${FF_BUILD_DIR}/lib/pkgconfig:$PKG_CONFIG_PATH" && \
        mkdir lib && \
        cp -a ${APP_DIR}/lib .; \
    fi && \
    acap-build . && \
    if [ "${INSTALL}" = "YES" ]; then eap-install.sh ${IP_ADDR} ${PASSWORD} install; fi && \
    if [ "${START}" = "YES" ]; then eap-install.sh start; fi

#-------------------------------------------------------------------------------
# Create output directory, we copy files from eap to host
#-------------------------------------------------------------------------------
RUN mkdir /opt/eap
RUN mv *.eap /opt/eap
RUN cd /opt/eap && \
    for file in *.eap; do \
        mv "$file" "${file%.eap}_sdk_${VERSION}.eap"; \
    done

