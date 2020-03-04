package payload

import (
	"fmt"
	"time"
)

// JSONTime with 'yyyy-MM-dd HH:mm:ss' format
type JSONTime time.Time

// JSONDate with 'yyyy-MM-dd' format
type JSONDate time.Time

// MarshalJSON time with 'yyyy-MM-dd HH:mm:ss' format
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// MarshalJSON time with 'yyyy-MM-dd HH:mm:ss' format
func (t JSONDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}
