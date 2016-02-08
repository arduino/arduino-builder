// +build api

package main

import (
	"io/ioutil"
	"path/filepath"

	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"github.com/gin-gonic/gin"
)

func listen(context map[string]interface{}) {
	router := gin.Default()

	router.POST("/compile", func(c *gin.Context) {
		logger := i18n.WebLogger{}
		logger.Init()
		context[constants.CTX_LOGGER] = &logger

		var json struct {
			Sketch types.Sketch `json:"sketch"`
			Fqbn   string       `json:"fqbn"`
		}

		err := c.BindJSON(&json)

		if err != nil {
			c.JSON(400, gin.H{"error": "Malformed JSON", "message": err.Error()})
			return
		}

		if json.Fqbn == "" {
			c.JSON(400, gin.H{"error": "Malformed JSON", "message": "Missing fqbn property"})
			return
		}

		context[constants.CTX_SKETCH] = &json.Sketch
		context[constants.CTX_FQBN] = json.Fqbn

		err = builder.RunBuilder(context)

		if err != nil {
			c.JSON(500, gin.H{"out": logger.Out(), "error": err.Error()})
			return
		}

		binaries := struct {
			Elf []byte `json:"elf,omitempty"`
			Bin []byte `json:"bin,omitempty"`
			Hex []byte `json:"hex,omitempty"`
		}{}

		elfPath := filepath.Join(*buildPathFlag, json.Sketch.MainFile.Name+".elf")
		binaries.Elf, _ = ioutil.ReadFile(elfPath)

		binPath := filepath.Join(*buildPathFlag, json.Sketch.MainFile.Name+".bin")
		binaries.Bin, _ = ioutil.ReadFile(binPath)

		hexPath := filepath.Join(*buildPathFlag, json.Sketch.MainFile.Name+".hex")
		binaries.Hex, _ = ioutil.ReadFile(hexPath)

		c.JSON(200, gin.H{"out": logger.Out(), "binaries": binaries})
	})

	router.Run(":" + *listenFlag)
}
