# GoBlast

Demonstrates a way to run periodic tasks in goroutines for a limited time.

## How to use it

`goblast` can be invoked with all configuration from the command line
or provided via a file. For help on how to run it use `goblast -h`.

## Examples

Below shows an example that *blasts* a URL with all defaults; HTTP method `GET`, 60 seconds, 1 blaster and 10 requests/s:

```sh
$ goblast --url https://example.host.com/path --duration 10
...
```

A more useful example would be one that shows more available options:

```sh
$ goblast --url https://example.host.com/path \
    --method post \ # HTTP method
    --rate 100 \    # req/s
    --duration 30 \ # 30 seconds
    --num 20 \      # 20 blasters
    --header "Authentication: Bearer auth-token" \ # More headers can be added with the same flag
    --body '{"name":"shenanigan"}' # HTTP request body
...
```

`goblast` also support configuration files written in YAML. See [docs/blast.yaml](./docs/blast.yaml) for the specification.

Some configuration in the file can be overrided by flags from the command line.

**WARN!** Use this program responsibly.

