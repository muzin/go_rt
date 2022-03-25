package colorstr

import (
	"fmt"
	"strconv"
	"testing"
)

func TestColorString(t *testing.T) {
	t.Run("Test Color String ", func(t *testing.T) {

		t.Logf(Black("Black"))
		t.Logf(Red("Red"))
		t.Logf(Yellow("Yellow"))
		t.Logf(Green("Green"))
		t.Logf(Cyan("Cyan"))
		t.Logf(Blue("Blue"))
		t.Logf(Purple("Purple"))
		t.Logf(White("White"))

		t.Logf(LightBlack("Light Black"))
		t.Logf(LightRed("Light Red"))
		t.Logf(LightYellow("Light Yellow"))
		t.Logf(LightGreen("Light Green"))
		t.Logf(LightCyan("Light Cyan"))
		t.Logf(LightBlue("Light Blue"))
		t.Logf(LightPurple("Light Purple"))
		t.Logf(LightWhite("Light White"))

		t.Logf("successfully")
	})
}

func TestSannerColorRange(t *testing.T) {
	t.Run("Test Sanner Color Range", func(t *testing.T) {

		var str string = "Hello World."
		var color = 30
		for i := color; i < 115; i++ {
			sprintf := fmt.Sprintf("\x1b[0;%dm%s%s\x1b[0m", i, str, strconv.Itoa(i))
			fmt.Println(sprintf)
		}

	})
}
