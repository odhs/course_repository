# Install Golang in Ubuntu

Download Golang

```sh
curl -OL <https://go.dev/dl/go1.23.0.linux-amd64.tar.gz>
sudo tar -C /usr/local -xvf go1.23.0.linux-amd64.tar.gz
```

Put on the end of the file

```sh
export PATH=$PATH:/usr/local/go/bin
```

Run for the terminal to recognize the changes:

```sh
source ~/.profile
```

Check if it is working

```sh
go version
```
