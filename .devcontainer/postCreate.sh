#!/bin/bash

conda init
go install github.com/magefile/mage@latest
cd native/extern/vcpkg && . bootstrap-vcpkg.sh && cd /workspaces/v2k

/workspaces/v2k/native/extern/vcpkg/vcpkg install protobuf
sudo ln -s /workspaces/v2k/native/extern/vcpkg/packages/protobuf_x64-linux/tools/protobuf/protoc-* /usr/local/bin/protoc

touch $HOME/.bash_profile
cat .devcontainer/.profile >> $HOME/.bash_profile
