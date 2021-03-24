package colorstr

import "fmt"

// 前景色   背景色
//
// 30  	40	  黑色
//
// 31  	41	  红色
//
// 32  	42	  绿色
//
// 33  	43    黄色
//
// 34  	44    蓝色
//
// 35  	45 	  紫色
//
// 36  	46 	  青色
//
// 37  	47	  白色

const (
	// 暗色
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite

	// 亮色
	textLightBlack  = 90
	textLightRed    = 91
	textLightGreen  = 92
	textLightYellow = 93
	textLightBlue   = 94
	textLightPurple = 95
	textLightCyan   = 96
	textLightWhite  = 97
)

func Black(str string) string {
	return textColor(textBlack, str)
}
func Red(str string) string {
	return textColor(textRed, str)
}
func Yellow(str string) string {
	return textColor(textYellow, str)
}
func Green(str string) string {
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	return textColor(textCyan, str)
}
func Blue(str string) string {
	return textColor(textBlue, str)
}
func Purple(str string) string {
	return textColor(textPurple, str)
}
func White(str string) string {
	return textColor(textWhite, str)
}

func LightBlack(str string) string {
	return textColor(textLightBlack, str)
}
func LightRed(str string) string {
	return textColor(textLightRed, str)
}
func LightYellow(str string) string {
	return textColor(textLightYellow, str)
}
func LightGreen(str string) string {
	return textColor(textLightGreen, str)
}
func LightCyan(str string) string {
	return textColor(textLightCyan, str)
}
func LightBlue(str string) string {
	return textColor(textLightBlue, str)
}
func LightPurple(str string) string {
	return textColor(textLightPurple, str)
}
func LightWhite(str string) string {
	return textColor(textLightWhite, str)
}

// 文字上色
func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}
