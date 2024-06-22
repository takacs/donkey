## How to load cards from Anki Export

Load functionality is currently supporting only one type of Anki export.

To copy your deck from Anki export with the following steps:

Click on the ⚙️ next to the deck you want to transfer and select `Export`
<img width="664" alt="Screenshot 2024-06-22 at 19 30 32" src="https://github.com/takacs/donkey/assets/44911031/37387b06-367b-4282-a27b-7504aed7ddc8">

Use Export format `Cards in Plain Text (.txt)` and uncheck `Include HTML and media references`.
<img width="661" alt="Screenshot 2024-06-22 at 19 31 01" src="https://github.com/takacs/donkey/assets/44911031/cc2d79a8-d43e-4404-8c91-aea28cae36e1">

Once exported run the command
```
donkey load <path_to_export>
```
