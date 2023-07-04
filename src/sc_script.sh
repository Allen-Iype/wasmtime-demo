rustup target add wasm32-wasi
rustc lib.rs --target wasm32-wasi
#wasmtime --dir=. --dir=/tmp hello.wasm arp.txt /tmp/somewhere.txt
wasmtime --dir=. lib.wasm arp.txt 
#cat /tmp/somewhere.txt