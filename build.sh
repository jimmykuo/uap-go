#!/usr/bin/env bash

cd uap-core
git checkout v0.10.0
cd ..

# Strip out empty lines and comments for conciseness:
# yaml=`cat uap-core/regexes.yaml | sed '/\s*#/d' | sed '/^\s*$/d'`

# insert constom.yaml
src="uap-core/regexes.yaml"
custom="custom.yaml"
match=`grep -n '# Android General Device Matching' $src|cut -d':' -f1`
# echo $match
insert_point=`expr $match - 2`
echo "insert_point:$insert_point"

# Strip out empty lines and comments
# SonyEricsson is equal to Sony
yaml=`sed "$insert_point r $custom" $src | sed '/\s*#/d' | sed '/^\s*$/d' | sed 's/brand_replacement: '\''SonyEricsson'\''/brand_replacement: '\''Sony'\''/g'`

# Build and format a Go file including our sources:
echo "package uaparser
var definitionYaml = []byte(\`$yaml\`)" | gofmt > uaparser/yaml.go

# make sure uaparser is ok
cd uaparser
go build
