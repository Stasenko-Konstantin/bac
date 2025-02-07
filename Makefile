# android ndk r24

build:
	go mod tidy
	fyne package -os android -appID bac.mobile -icon logo.png -name bac

debug:
	adb -d install -r bac.apk

install: build debug