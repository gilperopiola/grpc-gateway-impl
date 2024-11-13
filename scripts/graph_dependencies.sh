#!/bin/bash

# Internal
echo "digraph dependencies {" > ./docs/dependencies_internal.dot
go list -f '{{.ImportPath}} {{join .Imports " "}}' ./... | awk ' 
{
  for (i = 2; i <= NF; i++) {
    if ($1 !~ /gilperopiola/ || $i !~ /gilperopiola/ || \
        $1 ~ /\/etc(\/|$)/ || $i ~ /\/etc(\/|$)/) {
        continue;
    }
    print "\"" $1 "\" -> \"" $i "\"";
  }
}' >> ./docs/dependencies_internal.dot
echo "}" >> ./docs/dependencies_internal.dot
dot -Tpng ./docs/dependencies_internal.dot -o ./docs/dependencies_internal.png

## External
# echo "digraph dependencies {" > ./docs/dependencies_external.dot
# go mod graph | awk '{print "\"" $1 "\" -> \"" $2 "\";"}' >> ./docs/dependencies_external.dot
# echo "}" >> ./docs/dependencies_external.dot
# dot -Tpng ./docs/dependencies_external.dot -o ./docs/dependencies_external.png


##
##  if ($1 !~ /gilperopiola/ || $i !~ /gilperopiola/ || \
##      $1 ~ /\/etc(\/|$)/ || $i ~ /\/etc(\/|$)/ || \
##      $1 ~ /\/app\/clients(\/|$)/ || $i ~ /\/app\/clients(\/|$)/) {
##      continue;
##  }
##