# Flatten

Flattens the object - it'll return an object one level deep,
regardless of how nested the original object was.

## Install

```bash
go install github.com/josestg/go/flatten/cmd/flatten
```
make sure `GO111MODULE=on`.

## Examples

1. From JSON file.
```shell
# nested object.
cat nested.json | jq
{
  "a": {
    "b": [
      {
        "c": "1",
        "d": "2"
      },
      {
        "e": "3",
        "f": "4"
      }
    ]
  }
}

# flatten object.
cat nested.json | flatten | jq
{
  "a.b.0.c": "1",
  "a.b.0.d": "2",
  "a.b.1.e": "3",
  "a.b.1.f": "4"
}
```

2. From YAML file.

```shell
# nested object.
cat nested.yaml
a:
  b:
  - c: '1'
    d: '2'
  - e: '3'
    f: '4'
    
# flatten object.
cat nested.json | flatten | jq
{
  "a.b.0.c": "1",
  "a.b.0.d": "2",
  "a.b.1.e": "3",
  "a.b.1.f": "4"
}
```


