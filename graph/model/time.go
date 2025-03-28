package model

import (
	"fmt"
	"io"
	"time"
)

// Time is a custom GraphQL scalar for timestamps
type Time time.Time

// UnmarshalGQL converts a GraphQL input string to Go's time.Time
func (t *Time) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("time must be a string")
	}

	parsedTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	*t = Time(parsedTime)
	return nil
}

// MarshalGQL converts Go's time.Time to a GraphQL output
func (t Time) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, time.Time(t).Format(time.RFC3339))
}

// ToTime converts model.Time to time.Time
func (t Time) ToTime() time.Time {
	return time.Time(t)
}

// ToModelTime converts time.Time to model.Time
func ToModelTime(t time.Time) Time {
	return Time(t)
}
