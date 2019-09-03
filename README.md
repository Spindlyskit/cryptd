# Cryptd

Cryptd is a set of cipher crackers written in go. This application alone does not provide any way of interacting with the crackers aside from an rpc server meaning a client application will be required to utilize this project. Eventually I will create a cli client for interacting with Cryptd.

# Setup
In order to compile cryptd you must have go installed.
In order to setup cryptd for development you first need to clone the repo using `git clone https://github.com/Spindlyskit/cryptd.git`.
Next cd into the newly cloned repo and run `go get` to install the dependencies.
Finally run `make` to run tests and create binaries for linux, macOS and windows.
To delete these binaries run `make clean`

