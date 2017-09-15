package jsonrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	builder "github.com/arduino/arduino-builder"
	"github.com/arduino/arduino-builder/types"
	"github.com/fsnotify/fsnotify"
	"github.com/osamingo/jsonrpc"
	"golang.org/x/net/context"
)

type (
	BuildHandler struct {
		watcher *fsnotify.Watcher
		ctx     *types.Context
	}
	BuildParams struct {
		HardwareFolders         string
		ToolsFolders            string
		BuiltInLibrariesFolders string
		OtherLibrariesFolders   string
		SketchLocation          string
		FQBN                    string
		ArduinoAPIVersion       string
		CustomBuildProperties   string
		Verbose                 bool
		BuildCachePath          string
		BuildPath               string
		WarningsLevel           string
	}
	BuildResult struct {
		Message string `json:"message"`
		Error   error
	}
)

type (
	WatchHandler struct {
		watcher *fsnotify.Watcher
	}
	WatchParams struct {
		Path string `json:"path"`
	}
	WatchResult struct {
		Message string `json:"message"`
	}
)

func (h *BuildHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	var p BuildParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		fmt.Println(err)
		return nil, err
	}

	h.ctx.HardwareFolders = strings.Split(p.HardwareFolders, ",")
	h.ctx.ToolsFolders = strings.Split(p.ToolsFolders, ",")
	h.ctx.BuiltInLibrariesFolders = strings.Split(p.BuiltInLibrariesFolders, ",")
	h.ctx.OtherLibrariesFolders = strings.Split(p.OtherLibrariesFolders, ",")
	h.ctx.SketchLocation = p.SketchLocation
	h.ctx.CustomBuildProperties = strings.Split(p.CustomBuildProperties, ",")
	h.ctx.ArduinoAPIVersion = p.ArduinoAPIVersion
	h.ctx.FQBN = p.FQBN
	h.ctx.Verbose = p.Verbose
	h.ctx.BuildCachePath = p.BuildCachePath
	h.ctx.BuildPath = p.BuildPath
	h.ctx.WarningsLevel = p.WarningsLevel

	err := builder.RunBuilder(h.ctx)
	if err != nil {
		return BuildResult{
			Message: h.ctx.GetLogger().Flush(),
			Error:   err,
		}, nil
	}

	return BuildResult{
		Message: h.ctx.GetLogger().Flush(),
	}, nil
}

func (h *WatchHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	var p WatchParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	err := h.watcher.Add(p.Path)
	if err != nil {
		return nil, jsonrpc.ErrInvalidParams()
	}
	return BuildResult{
		Message: "OK " + p.Path,
	}, nil
}

func startWatching(ctx *types.Context) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				ctx.CanUseCachedTools = false
				log.Println("event:", event)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	return watcher
}

func RegisterAndServeJsonRPC(ctx *types.Context) {

	watcher := startWatching(ctx)

	jsonrpc.RegisterMethod("Main.Build", &BuildHandler{watcher, ctx}, BuildParams{}, BuildResult{})
	jsonrpc.RegisterMethod("Main.AddWatchPath", &WatchHandler{watcher}, WatchParams{}, WatchResult{})

	http.HandleFunc("/jrpc", func(w http.ResponseWriter, r *http.Request) {
		jsonrpc.HandlerFunc(r.Context(), w, r)
	})
	http.HandleFunc("/jrpc/debug", jsonrpc.DebugHandlerFunc)
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalln(err)
	}
}
