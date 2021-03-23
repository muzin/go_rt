package colorstr

import (
	"fmt"
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

		t.Logf("successfully")
	})
}

func TestCustomColorString(t *testing.T) {
	t.Run("Test Custom Color String", func(t *testing.T) {

		var str string = "Hello World."

		var color = 30

		for i := color; i < 1000; i++ {
			sprintf := fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", i, str)
			fmt.Println(sprintf)
		}

	})
}
