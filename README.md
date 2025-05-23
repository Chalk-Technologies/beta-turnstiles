
## Downloading
`chmod +x install.sh`
`./install.sh`

## Building
`go install fyne.io/tools/cmd/fyne@latest`
`fyne package -os darwin`
`fyne package -os linux`
`fyne package -os windows`

## Path issues
`export GOPATH=$HOME/go`
`export PATH=$GOPATH/bin:$PATH`



