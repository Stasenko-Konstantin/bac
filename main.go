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
	players           int
	currPlayer        = 1
	regime            string
	mysteryWord       string
	mysteryWordLength int
)

const (
	singleGame  = "одиночная"
	localGame   = "локальная"
	networkGame = "сетевая"
)

func main() {
	a := app.New()
	w := a.NewWindow("bac")

	//choosePlayers := widget.NewSelect(
	//	[]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
	//	func(p string) {
	//		i, _ := strconv.Atoi(p)
	//		players = i
	//	})
	//choosePlayers.SetSelected("2")

	var (
		newGameB            *widget.Button
		startGameC          *fyne.Container
		currPlayerL         = widget.NewLabel("")
		mysteryWordsLengthL = widget.NewLabel("")
		newWordE            = widget.NewEntry()
		givenWords          []string
		givenWordsL         = widget.NewLabel("")
		okB                 = widget.NewButton("ok", func() {})
	)

	_ = givenWords

	okB.OnTapped = func() {
		defer func() {
			newWordE.SetText("")
		}()

		if regime != singleGame {
			currPlayerL.SetText(fmt.Sprintf("текущий игрок: %d", currPlayer))
		}

		word := newWordE.Text
		wordLength := utf8.RuneCountInString(word)
		if wordLength != mysteryWordLength {
			text := fmt.Sprintf(
				"длина введенного слова\nне соответствует длине отгадываемого:\n%d и %d соответственно",
				wordLength,
				mysteryWordLength,
			)
			dialog.ShowCustomConfirm(
				text,
				"да",
				"тоже да",
				widget.NewLabel("я буду внимательней"),
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
			givenWordsL.SetText(fmt.Sprintf("%s - %dб, %dк\n%s", word, bulls, cows, givenWordsL.Text))
		}

		if word == mysteryWord {
			text := "победа!"
			if regime != singleGame {
				text = "" // TODO
			}
			dialog.ShowCustomConfirm("победа!", "еще!", "хватит...", widget.NewLabel(text), func(b bool) {
				if !b {
					os.Exit(1)
				}
				currPlayerL.SetText("")
				mysteryWordsLengthL.SetText("")
				givenWordsL.SetText("")
				startGameC.Hide()
				newGameB.Show()
			}, w)
			return
		}

		if regime != singleGame {
			if currPlayer == players {
				currPlayer = 1
			} else {
				currPlayer++
			}
		}
	}

	chooseRegime := widget.NewRadioGroup([]string{singleGame, localGame, networkGame}, func(p string) {
		regime = p
	})

	newGameB = widget.NewButton("Новая игра", func() {
		dialog.ShowCustomConfirm("Выберите режим игры:", "ok", "не хочу", chooseRegime, func(b bool) {
			if !b {
				startGameC.Hide()
				newGameB.Show()
			}

			if regime == singleGame {
				chooseWordLength := widget.NewSelect(
					[]string{"3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
					func(wl string) {
						i, _ := strconv.Atoi(wl)
						words := dict[i]
						mysteryWord = words[rand.Intn(len(words))]
						mysteryWordLength = utf8.RuneCountInString(mysteryWord)
						mysteryWordsLengthL.SetText(fmt.Sprintf("длина отгадываемого слова: %d", mysteryWordLength))
						currPlayerL.SetText(fmt.Sprintf("текущий игрок: %d", currPlayer))
						if regime == singleGame {
							currPlayerL.SetText("")
						}
						fmt.Println(mysteryWord)
					})
				chooseWordLength.SetSelected("3")
				dialog.ShowCustom("Выберите длинну угадываемого слова", "ok", chooseWordLength, w)
			}
		}, w)
		newGameB.Hide()
		startGameC.Show()
		currPlayer = 1
	})

	startGameC = container.NewGridWithRows(2,
		container.NewVBox(
			currPlayerL,
			mysteryWordsLengthL,
			newWordE,
			okB,
		),
		container.NewVScroll(givenWordsL),
	)

	startGameC.Hide()

	content := container.NewVBox(
		widget.NewLabel("\n\n\n\n"),
		newGameB,
		startGameC,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
