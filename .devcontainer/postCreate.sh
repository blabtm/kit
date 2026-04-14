#!/bin/bash

conda init
task rpk:profile:update

cd /workspaces/v2k/native/extern/vcpkg && . bootstrap-vcpkg.sh && cd /workspaces/v2k

/workspaces/v2k/native/extern/vcpkg/vcpkg install protobuf
sudo ln -sf                                                                               \
  /workspaces/v2k/native/extern/vcpkg/packages/protobuf_x64-linux/tools/protobuf/protoc-* \
  /usr/local/bin/protoc

cat .devcontainer/.profile >> $HOME/.bash_profile

