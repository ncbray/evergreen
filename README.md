# Evergreen

## About
Evergreen is a sandbox for exploring source-to-source compiler construction with
Go.

## Quick Start

### Setup
Be sure to set GOPATH to this directory.  If using bash:

    export GOPATH=`pwd`

Development workflow is assisted by a binary called "workflow".  Evergreen
generates parts of itself.  You need to build workflow to get started.

    go install evergreen/cmd/workflow

### Development
Regenerate into src/generated and run tests:

    ./bin/workflow test

Files in src/generated are only used for testing and will not affect the main
build.  This will help sanity check your changes to *.dub files.  If things look
OK, regenerate into src/evergreen:

    ./bin/workflow build

This command will compile the current program to bin/regenerate before replacing
the generated source files.  This is a safety net that lets generated sources be
updated even when the source code is in an intermediate, uncompilable state.  If
you need to regenerate sources while the compiler is broken:

    ./bin/workflow sync

This development loop should help minimize situations where work in progress is
lost because changes to the parser break it and prevent it from regenerating
itself.  This is still a risk, however, so be careful.

## Background

### Layout
* /dub/ contains the DSL sources.
* /src/ contains the Go sources.
* /src/generated/ contains a temporary copy of the generated sources for testing.
* /output/ is a temporary directory for visualization and debugging information.

### Terminology
* __text__: source code.
* __tree__: abstract syntax tree.
* __flow__: a graph based IR.

An agressive source-to-source compiler will go all the way down and back up this
stack:

text => tree => flow => flow => tree => text

Simpler source-to-source compilers will skip many of these stages, trading away
the ability to do deeper transformations for simplicity:

text => tree => text
