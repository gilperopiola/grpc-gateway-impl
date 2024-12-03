#!/bin/bash

# Internal
echo "digraph dependencies {" > ./docs/deps_internal.dot
go list -f '{{.ImportPath}} {{join .Imports " "}}' ./... | awk ' 
{
  for (i = 2; i <= NF; i++) {
    if ($1 !~ /gilperopiola/ || $i !~ /gilperopiola/ || \
        $1 ~ /\/mongo(\/|$)/ || $i ~ /\/mongo(\/|$)/ || \
        $1 ~ /\/easyql(\/|$)/ || $i ~ /\/easyql(\/|$)/ || \
        $1 ~ /\/etc(\/|$)/ || $i ~ /\/etc(\/|$)/) {
        continue;
    }
    print "\"" $1 "\" -> \"" $i "\"";
  }
}' >> ./docs/deps_internal.dot
echo "}" >> ./docs/deps_internal.dot
dot -Tpng ./docs/deps_internal.dot -o ./docs/deps_internal.png

## External
# echo "digraph dependencies {" > ./docs/deps_external.dot
# go mod graph | awk '{print "\"" $1 "\" -> \"" $2 "\";"}' >> ./docs/deps_external.dot
# echo "}" >> ./docs/deps_external.dot
# dot -Tpng ./docs/deps_external.dot -o ./docs/deps_external.png

## Example of filtering
##  if ($1 !~ /gilperopiola/ || $i !~ /gilperopiola/ || \
##      $1 ~ /\/etc(\/|$)/ || $i ~ /\/etc(\/|$)/ || \
##      $1 ~ /\/app\/clients(\/|$)/ || $i ~ /\/app\/clients(\/|$)/) {
##      continue;
##  }
##