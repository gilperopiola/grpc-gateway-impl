package b

import "github.com/gilperopiola/grpc-gateway-impl/etc/tests/import_cycle_test/a"

var B = 'b'

var Q = a.A
