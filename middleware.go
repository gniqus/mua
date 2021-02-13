package mua

import "log"

func ExceptionHandler(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			c.EchoString("500 Internal Server Error")
		}
	}()
	c.Next()
}
