package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPCookieJoinActionRun(t *testing.T) {
	t.Skip("Disabled due to temporary workaround for a bug") // FIXME: re-enable this test case

	oldCookies := map[string]string{
		"session": "12211",
		"n":       "1",
	}

	newCookies := map[string]string{
		"n": "2",
		"x": "1.11",
	}

	updatedCookies := map[string]string{
		"session": "12211",
		"n":       "2",
		"x":       "1.11",
	}

	oldCookiesIn := NewDataPipe()
	newCookiesIn := NewDataPipe()
	updatedCookiesOut := NewDataPipe()

	oldCookiesIn.Add(oldCookies)
	newCookiesIn.Add(newCookies)

	action := NewHTTPCookieJoinAction()

	action.AddInput(HTTPCookieJoinActionInputOldCookies, oldCookiesIn)
	action.AddInput(HTTPCookieJoinActionInputNewCookies, newCookiesIn)
	action.AddOutput(HTTPCookieJoinActionOutputUpdatedCookies, updatedCookiesOut)

	err := action.Run()

	assert.Nil(t, err)

	gotUpdatedCookies, ok := updatedCookiesOut.Remove().(map[string]string)

	assert.True(t, ok)
	assert.Equal(t, updatedCookies, gotUpdatedCookies)
}
