package encrypt

import (
	"fmt"
	"testing"
)

func TestEnPwdCode(t *testing.T) {
	uidString := "abcdsjjliesjdles"

	uidString, err := EnPwdCode([]byte(uidString))
	if err != nil {
		panic(err)
	}

	fmt.Println(uidString)
}

func TestDePwdCode(t *testing.T) {
	code, err := DePwdCode("gAwd1+TxWnh+xEXqxK8Hv4PRmHp2jPISBkqNevFC4Nw=")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(code))
}
