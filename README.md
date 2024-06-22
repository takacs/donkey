<h1 align=center> donkey </h1>
<p align="center">(d-anki)<p>
<p align="center">An Anki-like learning app for your terminal.<p>
<p align="center">
  <img src="assets/donkey.svg" width=20%/>
</p>
<p align="center">
  <img src="https://github.com/takacs/donkey/actions/workflows/ci.yml/badge.svg">
</p>
<p align="center">
  <img src="https://github.com/takacs/donkey/assets/44911031/132c1e06-7d91-46cc-bcfd-b05f5d28815d">
</p>

### Installation
#### Go
```
go install github.com/takacs/donkey@latest
```
### Runnning donkey
Running the `donkey` command after installation reaches the cli. To start the actual app:
```
donkey launch
```

### CLI:
There's some functionality that can be achieved through the cli

- Commands:
	- `donkey add`
		- add card to db
		- flags:
			- -f front
			- -b back
			- -d deck (optional: default)
	- `donkey list`
		- list all donkey cards
	- `donkey delete`
		- delete donkey cards based on id
		- args
			- -id id of card to delete
	- `donkey where`
		- show where the db is located on local machine
	- `donkey launch`
		- launch the Donkey TUI
	- `donkey load`
		- load anki exported file
		- flags:
			- -d deck (optional: default)

### Acknowledgments
Built with [bubbletea](https://github.com/charmbracelet/bubbletea)
