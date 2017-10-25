#!/bin/sh

virtualenv venv
. venv/bin/activate
pip install -r requirements.txt

brew install glide || :
brew install go || :
glide update


