#!/bin/bash

touch $HOME/.bash_profile
conda init
conda env create -f .devcontainer/environment.yaml

echo "conda activate v2k" >> $HOME/.bash_profile

