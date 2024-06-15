<h3 align="center">Donkey</h3>
<h4 align="center">(d-anki)</h4>
<h4 align="center">Study in your terminal.<h4>
<p align="center">
  <img src="assets/donkey.svg" width=35%/>
</p>
<p align="center">
  <img src="https://github.com/takacs/donkey/actions/workflows/ci.yml/badge.svg">
</p>


## donkey cli:

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
		- launch the [[Donkey TUI]]
