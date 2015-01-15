super basic Golang markdown parser - murk it! 

Removes most markdown syntax, giving clean, white and bold ANSI escaped console text.
pipe output through 'less -R' to see bold and color.

To build or install, you'll need to have your Go workspace setup - 
full instructions at https://golang.org/doc/code.html

but basically it's
`mkdir $HOME/go`
`export GOPATH=$HOME/go`

Requires github.com/mgutz/ansi - to install that, do:

`go get github.com/mgutz/ansi`

Then to actually build or install murkdown, it's just:

`git clone https://github.com/sideb0ard/murkdown.git` 

then `go build` or `go install` it. (go build will leave the binary in your current dir, whereas go install puts it in your $GOPATH/bin)


