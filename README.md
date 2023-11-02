
# Wasmtime-Go and Rust Integration Demo

This project serves as a demonstration for learning how `wasmtime-go` works and testing the integration of WebAssembly (WASM) with Go using Rust for the logic implementation. 

## Overview

The Rust components of this project implement various logical operations which are then compiled into WebAssembly (WASM) modules. These modules are integrated and used within Go test cases via `wasmtime-go`.

## Project Structure

The project has been refactored to enhance testability and maintainability:
- Go code relevant to each Rust contract has been moved to the corresponding Rust project folder.
- This refactoring aids in the easy testing of individual contracts and serves as a quick reference.

## Compiling Rust to WASM

To compile the Rust code to a WASM module, use the following command:

```sh
cargo build --target wasm32-unknown-unknown
