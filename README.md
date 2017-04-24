# yq

YAML wrapper for the fantastic [`jq`](https://stedolan.github.io/jq/) tool.

    go get github.com/gomatic/yq

# usage

    yq [options...] filter -- [files...]
    
Notice the `--`. It's mandatory so `yq` can distinguish the files from the options and filter.
It also allows `yq` to determine that input is coming from stdin (when stdin is a terminal).
