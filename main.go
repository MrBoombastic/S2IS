package main

import (
	"fmt"
	"github.com/BOOMfinity/golog"
	"github.com/MrBoombastic/S2fS/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
)

var log = golog.New("S2fS")

func main() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/upload", upload)
	app.Post("/delete", del)
	app.Use("/serve", filesystem.New(filesystem.Config{
		Root:   http.Dir("./s2fs_data"),
		Browse: false,
	}))
	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		if !fiber.IsChild() {
			log.Info().Send("Listening on: http://%v:%v", listenData.Host, listenData.Port)
		}
		return nil
	})
	port := ":" + os.Getenv("S2FS_PORT")
	if port == ":" {
		port = ":3000"
	}
	log.Fatal().SendError(app.Listen(port))
}

func upload(c *fiber.Ctx) error {
	var request types.UploadRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error().SendError(err)
		return c.Status(fiber.StatusBadRequest).JSON(types.UploadResponse{Error: err.Error()})
	}
	if request.Extension == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.UploadResponse{Error: "extension is empty"})
	}
	if request.File == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.UploadResponse{Error: "file is empty"})
	}
	// generate random uuid v4
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().SendError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(types.UploadResponse{Error: err.Error()})
	}
	// save file from request
	filename := fmt.Sprintf("%v.%v", id, request.Extension)
	if err := os.WriteFile("./s2fs_data/"+filename, []byte(request.File), 0644); err != nil {
		log.Error().SendError(err)
		return err
	}
	return c.JSON(types.UploadResponse{
		Filename: filename,
	})
}

func del(c *fiber.Ctx) error {
	var request types.DeleteRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error().SendError(err)
		return c.Status(fiber.StatusBadRequest).JSON(types.DeleteResponse{Error: err.Error()})
	}
	request.Filename = strings.ReplaceAll(request.Filename, "../", "")
	request.Filename = strings.ReplaceAll(request.Filename, "./", "")
	request.Filename = strings.ReplaceAll(request.Filename, "/", "")
	filename := fmt.Sprintf("./s2fs_data/%v", request.Filename)
	if err := os.Remove(filename); err != nil {
		log.Error().SendError(err)
		return err
	}
	return c.JSON(types.DeleteResponse{Error: ""})
}
