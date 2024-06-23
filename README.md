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

### What is Anki?

Anki is an open-source flash card program aimed to help users memorize and recall. For more check out Anki's wikipedia: [Anki (software)](https://en.wikipedia.org/wiki/Anki_(software))

### Installation
#### Go
```
go install github.com/takacs/donkey@latest
```
### Running donkey
Running the `donkey` command after installation reaches the cli. To start the actual app:
```
donkey launch
```

### Card Review
Cards to be reviewed are provided by a slightly modified version of the [SuperMemo-2 algorithm](https://en.wikipedia.org/wiki/SuperMemo).

### CLI
There's some functionality that can be achieved through the cli.

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
	- docs on how to import: https://github.com/takacs/donkey/blob/main/docs/load_cards.md 
	- load anki exported file
	- flags:
		- -d deck (optional: default)

### Acknowledgments
Built with [bubbletea](https://github.com/charmbracelet/bubbletea)
