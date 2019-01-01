package mpq

// #cgo windows CFLAGS: -D_MPQ_WINDOWS
// #cgo windows LDFLAGS: -Lstormlib -lstorm -lwininet -lz -lbz2 -lstdc++
// #cgo linux CFLAGS: -D_MPQ_LINUX
// #cgo linux LDFLAGS: -L./stormlib/ -lstorm -lz -lbz2 -lstdc++
import "C"
