#!/bin/bash

conda init

go install github.com/magefile/mage@latest

cd native/extern/vcpkg && . bootstrap-vcpkg.sh && cd /workspaces/v2k
/workspaces/v2k/native/extern/vcpkg/vcpkg install protobuf
sudo ln -sf                                                                               \
  /workspaces/v2k/native/extern/vcpkg/packages/protobuf_x64-linux/tools/protobuf/protoc-* \
  /usr/local/bin/protoc

curl -LO https://github.com/redpanda-data/redpanda/releases/latest/download/rpk-linux-amd64.zip
sudo unzip rpk-linux-amd64.zip -d /usr/local/bin
rm rpk-linux-amd64.zip
task rpk:profile:create

touch $HOME/.bash_profile
cat .devcontainer/.profile >> $HOME/.bash_profile

