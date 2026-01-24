package utils

import "github.com/sony/sonyflake"

func IDGenerate() (uint64, error) {
	settings := sonyflake.Settings{}
	sk := sonyflake.NewSonyflake(settings)

	id, err := sk.NextID()
	return id, err
}
