package main

// Cover struct
type Cover struct {
	height int
	url    string
	width  int
}

// Height resolves height field of Cover
func (c *Cover) Height() int32 {
	return int32(c.height)
}

// Width resolves width field of Cover
func (c *Cover) Width() int32 {
	return int32(c.width)
}

// URL resolves url field of Cover
func (c *Cover) URL() string {
	return c.url
}
