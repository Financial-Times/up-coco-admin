up-coco-admin
=============

[![CircleCI](https://circleci.com/gh/Financial-Times/up-coco-admin.svg?style=svg)](https://circleci.com/gh/Financial-Times/up-coco-admin)

Simple Universal Publishing stack admin tool for dev/ops/test purposes.

Initially, this application only contains a simple etcd value dump endpoint.  More will come later.

Installation:
-------------

```
go get github.com/Financial-Times/up-coco-admin
```

Usage
-----
```
up-coco-admin --help
```


API
---

See the [api.md](./api.md) file for details. To test the API documentation is still valid, simply install DreddJS and run the following:

```
# N.B. this is run with every Circle CI build
dredd --config ./api/api.yml
```
