#!/bin/bash
# Run from the top of the source tree
# Requires the test app to be built


${TOP}/example/bin/server &
pid=$!
sleep 3
${TOP}/example/bin/client
ret=$?
kill -9 $!
exit $ret
