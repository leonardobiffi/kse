# Kubernetes Secret Encoding

![release](https://img.shields.io/github/v/release/leonardobiffi/kse)
![workflow](https://img.shields.io/github/actions/workflow/status/leonardobiffi/kse/release.yml)

CLI Tool that allow encode and decode Kubernetes Secrets

## Features

- preserve yaml file order
- support multiples secrets in the same file
- stdin, file or directory as target
- finds all secrets in a target directory

## Install

### Requirements

- Install [jq](https://stedolan.github.io/jq/download/)
- Install [curl](https://curl.se/download.html)

```sh
curl -fsSL https://raw.githubusercontent.com/leonardobiffi/kse/master/scripts/install.sh | sh
```

### Using Golang

```sh
go install github.com/leonardobiffi/kse@latest
```

## Examples

### Decode

Searching for secrets in directory or file path

```sh
kse decode -o -f ./k8s/secrets/

kse decode -o -f ./k8s/secrets/secret.yaml
```

Using Stdin

```sh
kubectl get secret mysecret -o yaml | kse decode
```

```sh
cat secret.yaml | kse decode
```

### Encode

Searching for secrets in directory or file path

```sh
kse encode -o -f ./k8s/secrets/

kse encode -o -f ./k8s/secrets/secret.yaml
```

Using Stdin

```sh
cat secret.yaml | kse encode
```

## Roadmap

- Add support for JSON secrets
- Add tests
- Add Dockerfile
- Refactor packages encode/decode
