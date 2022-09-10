# Coin

Coin is my Golang blockchain implementation.

## Initialize First Node

First set NODE_ID.
```sh
export NODE_ID=3000 // linux
set NODE_ID=3000 // macos
```

Create wallet.
```sh
coin wallet --create
```

Create blockchain
```sh
coin server --create --address 1PprhXRdQQB5LjY5FNTNjkRbuxwc1E3Fh
```

Create start initial node.
```sh
coin server --address 1PprhXRdQQB5LjY5FNTNjkRbuxwc1E3Fh
```

## Initialize More Nodes

First set other NODE_ID.
```sh
export NODE_ID=5000 // linux
set NODE_ID=5000 // macos
```

Create wallet.
```sh
coin wallet --create
```

Start node.
```sh
coin server --address 1PprhXRdQQB5LjY5FNTNjkRbuxwc1E3Fh
```
