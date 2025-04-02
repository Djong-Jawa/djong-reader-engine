package model

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
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

// TimeFormat defines the layout for the time representation
const TimeFormat = time.RFC3339

// MarshalTime serializes a time.Time to a GraphQL scalar
func MarshalTime(t time.Time) graphql.Marshaler {
    return graphql.WriterFunc(func(w io.Writer) {
        formattedTime := t.Format(TimeFormat)
        io.WriteString(w, fmt.Sprintf("%q", formattedTime))
    })
}

// UnmarshalTime deserializes a GraphQL scalar to a time.Time
func UnmarshalTime(v interface{}) (time.Time, error) {
    switch v := v.(type) {
    case string:
        return time.Parse(TimeFormat, v)
    case *string:
        if v == nil {
            return time.Time{}, nil
        }
        return time.Parse(TimeFormat, *v)
    default:
        return time.Time{}, fmt.Errorf("unexpected type %T for Time", v)
    }
}

func TimeFromTime(t time.Time) Time {
    return Time(t)
}

func TimeFromPtr(t *time.Time) *Time {
    if t == nil {
        return nil
    }
    mt := Time(*t)
    return &mt
}
