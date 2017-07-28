Photocache
==========

Caching proxy designed to retrieve photos from popular photo services and store
them locally.

Getting Started
---------------

**Deploying to Spine**

1. Sign up for a [Spine](https://spi.ne) account if you do not yet have one.
2. Follow the CLI setup instructions provided within the signup e-mail.
3. Run ```spine deploy```.
4. 

Developer Instructions
----------------------

Building
--------

1. Create a directory for your $GOPATH and set it, if one does not yet exist.

```
mkdir ~/go
export GOPATH=~/go
```

2. Create a namespaced symlink for the current project in your $GOPATH.

```
mkdir -p "$GOPATH/src/github.com/ssalevan"
ln -s "$(pwd)" "$GOPATH/src/github.com/ssalevan/photocache"
```

3. Install Glide:

```
# brew install glide

or for the adventurous:

curl https://glide.sh/get | sh
```

4. Install all dependencies:

```
# glide install
```

5. Give it a build:

```
# go build .
```