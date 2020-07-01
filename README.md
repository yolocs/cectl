# cectl

To make operating CloudEvents easier via CLI.

## Installation

Grab binary.

```bash
go install github.com/yolocs/cectl
```

To use pre-built alpine based container image ([Dockerfile](./build/Dockerfile-alpine))

```bash
docker run -it cshou/cectl /bin/sh
```

To build it in your own container image, refer to [Dockerfile](./build/Dockerfile-alpine).

## Try me

To trigger your bash script (or whatever command) from CloudEvents.

```bash
cectl listen -p 8080 --cmd 'printenv'
```

To send a CloudEvent.

```bash
cectl send \
  --source cectl.source \
  --type cectl.roundtrip \
  --subject hello \
  --data 'custom-data' \
  --extensions "foo=bar" \
  --target "http://127.0.0.1:8080"
```

CloudEvent attributes could also be specified with [env vars](./pkg/env/env.go).

`printenv` command should be triggered and print out all the envs.

```
CE_IN_ID=923dfc03-d3d6-4aca-854d-589c2f1bba7b
CE_IN_SOURCE=cectl.source
CE_IN_TYPE=cectl.roundtrip
CE_IN_SUBJECT=hello
CE_IN_TIME=1593623591
CE_IN_DATASCHEMA=
CE_IN_CONTENTTYPE=
CE_IN_DATA="custom-data"
CE_IN_EXT_FOO=bar
```

All CloudEvent attributes were passed in as env vars which could be used in your script.

## References

[CloudEvent Spec](https://github.com/cloudevents/spec)
[CE Go SDK](https://github.com/cloudevents/sdk-go)