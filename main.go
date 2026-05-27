package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/zstd"
)

var decoder *zstd.Decoder

func main() {
	var err error
	decoder, err = zstd.NewReader(nil)
	if err != nil {
		panic(err)
	}
	defer decoder.Close()

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/decompress", handleDecompress)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func handleDecompress(c *gin.Context) {
	const MAX_BODY = 50 * 1024 * 1024

	compressed, err := io.ReadAll(io.LimitReader(c.Request.Body, MAX_BODY))
	if err != nil {
		c.String(http.StatusBadRequest, "failed to read body")
		return
	}
	if len(compressed) == 0 {
		c.String(http.StatusBadRequest, "empty body")
		return
	}

	decompressed, err := decoder.DecodeAll(compressed, nil)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, "decompression failed: "+err.Error())
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", decompressed)
}
