@echo off
setlocal

set ROOT_DIR=D:/code/go/gc
set PB_DIR=%ROOT_DIR%/examples/protobuf
set DIR_SOURCE=%PB_DIR%/source
set GEN_DIR=%PB_DIR%/genProto.exe

cd %DIR_SOURCE%

for /d %%a in (*) do (
    pushd %%a
    for %%b in (*.proto) do (
        %GEN_DIR% -s %DIR_SOURCE%/%%a/ -o %DIR_OUT%/%%a/

        protoc --go_out=. %%b
        protoc --go-grpc_out=. %%b
    )
    popd
)

echo ok

cd %ROOT_DIR%
endlocal