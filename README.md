# yq

[![Build Status](https://travis-ci.org/gomatic/yq.svg?branch=master)](https://travis-ci.org/gomatic/yq)

YAML wrapper for the fantastic [`jq`](https://stedolan.github.io/jq/) tool.

    go get github.com/gomatic/yq

# usage

    yq [options...] filter -- [files...]
    
Notice the `--`. It's mandatory so `yq` can distinguish the files from the options and filter.
It also allows `yq` to determine that input is coming from stdin.

# issues

Since `jq` can output JSON that isn't well-formed, such output can't _easily_ be converted back to yaml.
In such cases, the JSON will be output.
