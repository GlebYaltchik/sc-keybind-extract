# sc-keybind-extract
Simple tool to extract all Star Citizen keybinds to a CVS file

# Usage

## Before You start
- download and install `unp4k` utils written by Peter Dolkens from https://github.com/dolkensp/unp4k
- locate game data file `Data.p4k` \
  you can find the complete file location in the **Star Citizen Laucher** menu \
  `Settings -> Games -> Star Citizen - LIVE`

## Extract and prepare data

Let's assume what the game data located at `C:\Program Files\Star Citizen\LIVE`
```
unp4k.exe C:\Program Files\Star Citizen\LIVE\Data.p4k *.ini
unp4k.exe C:\Program Files\Star Citizen\LIVE\Data.p4k *.xml
```
That will create folders `Data` and `Engine` in the current directory.

Now we have to unpack compressed xml data.
```
unforge.exe Data\Libs\Config\defaultProfile.xml
unforge.exe Data\Libs\Config\keybinding_localization.xml 
```

## Build CSV export for keybinding

```
sc-keybind-extract --profile Data\Libs\Config\defaultProfile.xml --localization Data\Localization\english\global.ini > kbd.csv
```

Additionally, you can add information about keymap changes comparing to previous version of the game.
To do that you have to provide the previous version of the game data as follows:

```
sc-keybind-extract --profile profile-3.23.1.xml --prev-profile profile-3.23.0.xml --localization Data\Localization\english\global.ini > kbd.csv
```

ENJOY!

# Install

Build from scratch (go version 1.22 or more required)
```
go install -v github.com/GlebYaltchik/sc-keybind-extract@latest 
```