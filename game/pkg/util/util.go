// Package util 提供与业务无关的通用工具函数。
package util

// MaskPhone 将手机号中间四位替换为 *，如 "13800001111" → "138****1111"。
func MaskPhone(p string) string {
	if len(p) >= 11 {
		return p[:3] + "****" + p[len(p)-4:]
	}
	return p
}
