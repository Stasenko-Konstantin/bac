# bac
"Bulls and cows” - game on android in which you have to guess words <br>
supports only Russian

### build:
    fyne package -os android -appID bac.mobile -icon logo.png -name bac

#### or:
    make build

### debug:
    adb -d install -r bac.apk

#### or:
    make debug

### install:
    make install

<br>

REQUIREMENTS: Golang, Android NDK, fyne, Android itself