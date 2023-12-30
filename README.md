[![Test](https://github.com/xryuseix/ctfcli-unit-test/actions/workflows/test.yaml/badge.svg)](https://github.com/xryuseix/ctfcli-unit-test/actions/workflows/test.yaml)

# Unit Test Tool for [CTFd/ctfcli](https://github.com/CTFd/ctfcli)

CI tool to test ctfcli flag formats.

<div align="center">
  <img src="./demo.png" width="80%">
</div>

## Usage

### Prepreation

prepare directory like **[./example](https://github.com/xryuseix/ctfcli-unit-test/tree/main/example)** ([ctfcli](https://github.com/CTFd/ctfcli) format).

```txt
web/
  ├chall1/
          ├── challenge.yml
          ├── flag.txt
  ├chall2/
          ├── challenge.yml
          ├── flag.txt
misc/
  ├ ...
```

### Use with GitHub Actions

```yaml
name: Unit Test for CTFd/ctfcli

on:
  pull_request:
    branches: 
      - main

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Unit Test for CTFd/ctfcli
        uses: xryuseix/ctfcli-unit-test@v1.0.0
        with:
          target_directory: example
```

### Use with Command Line

```bash
# try to run
make run
# for production
# change INPUT_TARGET_DIRECTORY
make build && INPUT_TARGET_DIRECTORY="example" ./out
```
