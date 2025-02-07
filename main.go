package main

// TODO особая благодарность bark-arf за лучший UX

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"os"
	"strconv"
	"unicode/utf8"
)

var (
	mysteryWord       string
	mysteryWordLength int
)

func main() {
	a := app.New()
	w := a.NewWindow("bac")

	var (
		newGameB            *widget.Button
		newGameC            *fyne.Container
		rulesL              = widget.NewLabel(rules)
		startGameC          *fyne.Container
		mysteryWordsLengthL = widget.NewLabel("")
		newWordE            = widget.NewEntry()
		givenWords          = map[string]struct{}{}
		givenWordsL         = widget.NewLabel("")
		okB                 = widget.NewButton("ok", func() {})
		againB              = widget.NewButton("заново", func() {})
	)

	_ = givenWords

	okB.OnTapped = func() {
		defer func() {
			newWordE.SetText("")
		}()

		word := newWordE.Text
		wordLength := utf8.RuneCountInString(word)
		if wordLength != mysteryWordLength {
			text := fmt.Sprintf(
				"длина введенного слова\nне соответствует длине отгадываемого:\n%d и %d соответственно.",
				wordLength,
				mysteryWordLength,
			)
			dialog.ShowCustomConfirm(
				text,
				"да",
				"тоже да",
				widget.NewLabel("я буду внимательней!"),
				func(b bool) {
					newGameB.SetText("")
				},
				w,
			)
		} else {
			var cows, bulls int
			mysteryRunes := []rune(mysteryWord)
			wordRunes := []rune(word)
			checkedMystery := make([]bool, len(mysteryRunes))
			checkedWord := make([]bool, len(wordRunes))
			for i := range wordRunes {
				if wordRunes[i] == mysteryRunes[i] {
					bulls++
					checkedMystery[i] = true
					checkedWord[i] = true
				}
			}
			for i := range wordRunes {
				if checkedWord[i] {
					continue
				}
				for j := range mysteryRunes {
					if !checkedMystery[j] && wordRunes[i] == mysteryRunes[j] {
						cows++
						checkedMystery[j] = true
						break
					}
				}
			}
			if _, ok := givenWords[word]; !ok {
				givenWordsL.SetText(fmt.Sprintf("%s - %dб, %dк\n%s", word, bulls, cows, givenWordsL.Text))
			}
			givenWords[word] = struct{}{}
		}

		if word == mysteryWord {
			text := "ты угадал!"
			dialog.ShowCustomConfirm("победа!", "еще!", "хватит...", widget.NewLabel(text), func(b bool) {
				if !b {
					os.Exit(1)
				}
				mysteryWordsLengthL.SetText("")
				givenWordsL.SetText("")
				startGameC.Hide()
				newGameC.Show()
			}, w)
			return
		}
	}

	againB.OnTapped = func() {
		mysteryWordsLengthL.SetText("")
		givenWordsL.SetText("")
		startGameC.Hide()
		newGameC.Show()
	}

	newGameB = widget.NewButton("Новая игра", func() {
		chooseWordLength := widget.NewSelect(
			[]string{"3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
			func(wl string) {
				i, _ := strconv.Atoi(wl)
				words := dict[i]
				mysteryWord = words[rand.Intn(len(words))]
				mysteryWordLength = utf8.RuneCountInString(mysteryWord)
				mysteryWordsLengthL.SetText(fmt.Sprintf("длина отгадываемого слова: %d", mysteryWordLength))
				fmt.Println(mysteryWord)
			})
		chooseWordLength.SetSelected("3")
		dialog.ShowCustom("Выберите длинну угадываемого слова", "ok", chooseWordLength, w)

		newGameC.Hide()
		startGameC.Show()
	})

	newGameC = container.NewVBox(
		newGameB,
		rulesL,
	)

	startGameC = container.NewGridWithRows(2,
		container.NewVBox(
			againB,
			mysteryWordsLengthL,
			newWordE,
			okB,
		),
		container.NewVScroll(givenWordsL),
	)

	startGameC.Hide()

	content := container.NewVBox(
		newGameC,
		startGameC,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
