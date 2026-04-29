FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

COPY .profile /home/vscode/.bash_profile

RUN apt update && apt install -y \
  pkg-config

ENV VCPKG_ROOT="/opt/vcpkg"
ENV VCPKG_BINARY_SOURCES="clear;files,/opt/vcpkg_cache,readwrite"
ENV VCPKG_INSTALLED_ROOT="/opt/vcpkg_installed"
ENV CMAKE_TOOLCHAIN_FILE="/opt/vcpkg/scripts/buildsystems/vcpkg.cmake"

RUN git clone https://github.com/microsoft/vcpkg.git ${VCPKG_ROOT} \
    && ${VCPKG_ROOT}/bootstrap-vcpkg.sh

WORKDIR /tmp/vcpkg_manifest

COPY native/vcpkg.json ./
COPY native/vcpkg-configuration.json ./

RUN ${VCPKG_ROOT}/vcpkg install               \
    --x-manifest-root=/tmp/vcpkg_manifest     \
    --x-install-root=${VCPKG_INSTALLED_ROOT}

RUN ln -sf "${VCPKG_ROOT}/packages/protobuf_x64-linux/tools/protobuf/protoc" \
           "/usr/local/bin/protoc"

RUN chown -R vscode:vscode /opt/vcpkg_installed

