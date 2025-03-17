package main

import (
	"testing"
	"time"

	"github.com/markponce/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {

	//// Initial Test

	// // Initialize a new time.Time object and pass it to the humanDate function.
	// tm := time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC)
	// hd := humanDate(tm)

	// // Check that the output from the humanDate function is in the format we
	// // expect. If it isn't what we expect, use the t.Errorf() function to
	// // indicate that the test has failed and log the expected and actual
	// // values.
	// if hd != "17 Mar 2024 at 10:15" {
	// 	t.Errorf("got %q; want %q", hd, "17 Mar 2024 at 10:15")
	// }

	// Table of test cases
	// Slice of Anonymous struct
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			// if hd != tt.want {
			// 	t.Errorf("got %q, want %q", hd, tt.want)
			// }
			assert.Equal(t, hd, tt.want)
		})
	}

	// Exampel of sub-tests without a table of test cases
	// t.Run("Example sub-test 1", func(t *testing.T) {
	//     // Do a test.
	// })

	// t.Run("Example sub-test 2", func(t *testing.T) {
	//     // Do another test.
	// })

	// t.Run("Example sub-test 3", func(t *testing.T) {
	//     // And another...
	// })
}
