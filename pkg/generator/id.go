package generator

import "github.com/sony/sonyflake"

// IDGenerate 唯一ID生成器
func IDGenerate() (uint64, error) {
	settings := sonyflake.Settings{}
	sf := sonyflake.NewSonyflake(settings)

	id, err := sf.NextID()
	return id, err
}