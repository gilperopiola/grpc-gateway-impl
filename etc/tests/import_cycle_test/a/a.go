package a

import "github.com/gilperopiola/grpc-gateway-impl/etc/tests/import_cycle_test/d"

var A = 'a'

var Q = d.B
var K = d.C
