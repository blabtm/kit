FROM ghcr.io/blabtm/v2k.dev:latest AS build

WORKDIR /workspace

COPY CMakeLists.txt .
COPY CMakePresets.json .
COPY vcpkg.json .
COPY vcpkg-configuration.json .

COPY extern ./extern
COPY em-es.cpp ./em-es.cpp

RUN cmake --preset dev
RUN cmake --build build --target svc-em-es

FROM ubuntu:24.04

COPY --from=build                       \
  /workspace/build/em-es.cpp/svc-em-es  \
  /usr/local/bin/svc-em-es

CMD ["/usr/local/bin/svc-em-es"]
