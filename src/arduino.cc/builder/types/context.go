package types

import "strings"
import "arduino.cc/builder/i18n"
import "arduino.cc/builder/props"

// Context structure
type Context struct {
	// Build options
	HardwareFolders         []string
	ToolsFolders            []string
	LibrariesFolders        []string
	BuiltInLibrariesFolders []string
	OtherLibrariesFolders   []string
	SketchLocation          string
	ArduinoAPIVersion       string
	FQBN                    string

	BuildPath          string
	SketchBuildPath    string
	CoreBuildPath      string
	LibrariesBuildPath string
	PreprocPath        string

	WarningsLevel string

	// Libraries handling
	Includes          []string
	Libraries         []*Library
	HeaderToLibraries map[string][]*Library

	// C++ Parsing
	CTagsOutput                 string
	CTagsTargetFile             string
	CTagsOfSource               []*CTag
	CTagsOfPreprocessedSource   []*CTag
	CTagsCollected              []*CTag
	PrototypesSection           string
	PrototypesLineWhereToInsert int
	Prototypes                  []*Prototype

	// Verbosity settings
	Verbose           bool
	DebugPreprocessor bool

	// Contents of a custom build properties file (line by line)
	CustomBuildProperties []string

	// Logging
	logger     i18n.Logger
	DebugLevel int

	// For now it is used in conjunction with the old map[string]string, but
	// it will be slowly populated with all the fields currently used in the
	// map in the next commits.
	// When the migration will be completed the old map will be removed.
}

func (ctx *Context) ExtractBuildOptions() props.PropertiesMap {
	opts := make(props.PropertiesMap)
	opts["hardwareFolders"] = strings.Join(ctx.HardwareFolders, ",")
	opts["toolsFolders"] = strings.Join(ctx.ToolsFolders, ",")
	opts["builtInLibrariesFolders"] = strings.Join(ctx.BuiltInLibrariesFolders, ",")
	opts["otherLibrariesFolders"] = strings.Join(ctx.OtherLibrariesFolders, ",")
	opts["sketchLocation"] = ctx.SketchLocation
	opts["fqbn"] = ctx.FQBN
	opts["runtime.ide.version"] = ctx.ArduinoAPIVersion
	opts["customBuildProperties"] = strings.Join(ctx.CustomBuildProperties, ",")
	return opts
}

func (ctx *Context) InjectBuildOptions(opts props.PropertiesMap) {
	ctx.HardwareFolders = strings.Split(opts["hardwareFolders"], ",")
	ctx.ToolsFolders = strings.Split(opts["toolsFolders"], ",")
	ctx.BuiltInLibrariesFolders = strings.Split(opts["builtInLibrariesFolders"], ",")
	ctx.OtherLibrariesFolders = strings.Split(opts["otherLibrariesFolders"], ",")
	ctx.SketchLocation = opts["sketchLocation"]
	ctx.FQBN = opts["fqbn"]
	ctx.ArduinoAPIVersion = opts["runtime.ide.version"]
	ctx.CustomBuildProperties = strings.Split(opts["customBuildProperties"], ",")
}

func (ctx *Context) GetLogger() i18n.Logger {
	if ctx.logger == nil {
		return &i18n.HumanLogger{}
	}
	return ctx.logger
}

func (ctx *Context) SetLogger(l i18n.Logger) {
	ctx.logger = l
}
