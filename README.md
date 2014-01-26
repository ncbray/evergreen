About
=====
Evergreen is a sandbox for exploring compiler construction with Go.

Setup
=====
Be sure to set GOPATH to this directory.  If using bash:

  export GOPATH=`pwd`

Development
===========
Regenerate the parser:

  go run src/tools/regenerate/regenerate_main.go

This will generate the parser into src/generated. Sanity check the generated
parser:

  go test ./...

Now that it's known the generated parser passes tests, regenerate the parser in
place:

  go run src/tools/regenerate/regenerate_main.go -replace
  go fmt ./...

This development loop should help minimize situations where work in progress is
lost because changes to the parser break it and prevent it from regenerating
itself.
