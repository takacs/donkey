## How to load cards from Anki Export

Load functionality is currently supporting only one type of Anki export.

To copy your deck from Anki export with the following steps:

Click on the ⚙️ next to the deck you want to transfer and select `Export`
<img width="656" alt="Screenshot 2024-06-22 at 19 59 08" src="https://github.com/takacs/donkey/assets/44911031/8f33c7b7-2235-436e-b77c-28312952c4be">

Use Export format `Cards in Plain Text (.txt)` and uncheck `Include HTML and media references`.
<img width="661" alt="Screenshot 2024-06-22 at 19 31 01" src="https://github.com/takacs/donkey/assets/44911031/cc2d79a8-d43e-4404-8c91-aea28cae36e1">

Once exported run the command
```
donkey load <path_to_export>
```
