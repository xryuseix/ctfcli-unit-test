[![Test](https://github.com/xryuseix/ctfcli-unit-test/actions/workflows/test.yaml/badge.svg)](https://github.com/xryuseix/ctfcli-unit-test/actions/workflows/test.yaml)

# ctfcli Unit Test

Check your flag format for ctfcli challenge file using GitHub Actions.

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
config.yaml
```

config_file is optional. Create the following.

```yaml:config.yaml
# config.yaml
genre:
  - web
  - misc
  - osint
```

```yaml:challenge.yml
# challenge.yml
flags:
  - "flag{flag}"
  - {
      type: "regex",
      content: "flag{hoge.*}",
      data: "case_insensitive",
    }
```

```txt:flag.txt
# flag.txt
flag{flag} # this is comment
flag{Hoge123} # you can use regex or case_insensitive
!flag{flag} # expect this flag to be wrong(assert_failed).
flag{flag\#escaped} # you can escape `#`
```

### Use with GitHub Actions

```yaml
name: Flag Check

on:
  pull_request:
    branches: 
      - main

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Check Flag Format
        uses: xryuseix/ctfcli-unit-test@v2.0.1
        with:
          target_directory: example
          config_file: example/config.yaml
```

### Use with Command Line

```bash
# try to run
make run
# for production
# change INPUT_TARGET_DIRECTORY and INPUT_CONFIG_FILE
make build && INPUT_TARGET_DIRECTORY="example" INPUT_CONFIG_FILE="example/config.yaml" ./out
```
