About
=====
Evergreen is a sandbox for exploring compiler construction with Go.

Setup
=====
Be sure to set GOPATH to this directory.

Development
===========
Regenerate the parser:

  go run src/tools/regenerate/regenerate_main.go

Sanity check the generated parser:

  go test ./...

Copy the generated parser into place:

  go fmt ./...
  cp src/generated/dub/tree/generated_parser.go src/evergreen/dub/tree

Re-run the regenerate and testing steps.
