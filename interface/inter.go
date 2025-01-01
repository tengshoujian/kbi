type Base struct {
	b int
}
type Container struct {
	Base       //Container 是嵌入体结构
	c    sting // Base 是被嵌入体结构
}