package cryptocurrency

import (
	"strconv"
	"strings"
)

type Price float64

func (p Price) String() string {
	if p > 1 {
		return strconv.FormatFloat(float64(p), 'f', 2, 64)
	}
	return strconv.FormatFloat(float64(p), 'f', 10, 64)
}

type Metadata struct {
	USDT   Price  `json:"usdt_price"`
	Name   string `json:"name"`
	Change float64
	Vol    float64 `json:"vol"`
}

func (p Metadata) String() string {
	sb := &strings.Builder{}
	sb.WriteString(`Name: <b>`)
	sb.WriteString(strings.ToUpper(p.Name))
	sb.WriteString("</b>\n")

	sb.WriteString("USDT: <b>")
	sb.WriteString(p.USDT.String())
	sb.WriteString("</b>\n")

	if p.Change != 0 {
		sb.WriteString("Change 24h: ")
		if p.Change > 0 {
			sb.WriteString(`+`)
		}

		sb.WriteString(strconv.FormatFloat(p.Change, 'f', 2, 64))
		sb.WriteString("%\n")
	}

	if p.Vol != 0 {
		sb.WriteString("Value: ")
		sb.WriteString(strconv.FormatFloat(p.Vol, 'f', 2, 64))
		sb.WriteString("\n")
	}

	return sb.String() // strings.Replace(, ".", `\\.`, -1)

	//return str
}
