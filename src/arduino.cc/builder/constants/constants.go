/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 * Copyright 2015 Matthijs Kooijman
 */

package constants

const BOARD_PROPERTIES_MENU = "menu"
const BUILD_OPTIONS_FILE = "build.options.json"
const BUILD_PROPERTIES_ARCHIVE_FILE = "archive_file"
const BUILD_PROPERTIES_ARCHIVE_FILE_PATH = "archive_file_path"
const BUILD_PROPERTIES_ARCH_OVERRIDE_CHECK = "architecture.override_check"
const BUILD_PROPERTIES_BOOTLOADER_FILE = "bootloader.file"
const BUILD_PROPERTIES_BOOTLOADER_NOBLINK = "bootloader.noblink"
const BUILD_PROPERTIES_BUILD_ARCH = "build.arch"
const BUILD_PROPERTIES_BUILD_BOARD = "build.board"
const BUILD_PROPERTIES_BUILD_CORE = "build.core"
const BUILD_PROPERTIES_BUILD_CORE_PATH = "build.core.path"
const BUILD_PROPERTIES_BUILD_MCU = "build.mcu"
const BUILD_PROPERTIES_BUILD_PATH = "build.path"
const BUILD_PROPERTIES_BUILD_PROJECT_NAME = "build.project_name"
const BUILD_PROPERTIES_BUILD_SYSTEM_PATH = "build.system.path"
const BUILD_PROPERTIES_BUILD_VARIANT = "build.variant"
const BUILD_PROPERTIES_BUILD_VARIANT_PATH = "build.variant.path"
const BUILD_PROPERTIES_COMPILER_C_ELF_FLAGS = "compiler.c.elf.flags"
const BUILD_PROPERTIES_COMPILER_CPP_FLAGS = "compiler.cpp.flags"
const BUILD_PROPERTIES_COMPILER_PATH = "compiler.path"
const BUILD_PROPERTIES_COMPILER_WARNING_FLAGS = "compiler.warning_flags"
const BUILD_PROPERTIES_EXTRA_TIME_DST = "extra.time.dst"
const BUILD_PROPERTIES_EXTRA_TIME_LOCAL = "extra.time.local"
const BUILD_PROPERTIES_EXTRA_TIME_UTC = "extra.time.utc"
const BUILD_PROPERTIES_EXTRA_TIME_ZONE = "extra.time.zone"
const BUILD_PROPERTIES_INCLUDES = "includes"
const BUILD_PROPERTIES_OBJECT_FILE = "object_file"
const BUILD_PROPERTIES_OBJECT_FILES = "object_files"
const BUILD_PROPERTIES_PATTERN = "pattern"
const BUILD_PROPERTIES_PID = "pid"
const BUILD_PROPERTIES_PREPROCESSED_FILE_PATH = "preprocessed_file_path"
const BUILD_PROPERTIES_RUNTIME_HARDWARE_PATH = "runtime.hardware.path"
const BUILD_PROPERTIES_RUNTIME_OS = "runtime.os"
const BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH = "runtime.platform.path"
const BUILD_PROPERTIES_RUNTIME_TOOLS_PREFIX = "runtime.tools."
const BUILD_PROPERTIES_RUNTIME_TOOLS_SUFFIX = ".path"
const BUILD_PROPERTIES_RUNTIME_IDE_VERSION = "runtime.ide.version"
const BUILD_PROPERTIES_SOFTWARE = "software"
const BUILD_PROPERTIES_SOURCE_FILE = "source_file"
const BUILD_PROPERTIES_SOURCE_PATH = "build.source.path"
const BUILD_PROPERTIES_TOOLS_KEY = "tools"
const BUILD_PROPERTIES_VID = "vid"
const CTAGS = "ctags"
const EMPTY_STRING = ""
const FILE_BOARDS_LOCAL_TXT = "boards.local.txt"
const FILE_BOARDS_TXT = "boards.txt"
const FILE_BUILTIN_TOOLS_VERSIONS_TXT = "builtin_tools_versions.txt"
const FILE_CTAGS_TARGET = "ctags_target.cpp"
const FILE_CTAGS_TARGET_FOR_GCC_MINUS_E = "ctags_target_for_gcc_minus_e.cpp"
const FILE_GCC_PREPROC_TARGET = "gcc_preproc_target.cpp"
const FILE_PLATFORM_KEYS_REWRITE_TXT = "platform.keys.rewrite.txt"
const FILE_PLATFORM_LOCAL_TXT = "platform.local.txt"
const FILE_PLATFORM_TXT = "platform.txt"
const FILE_PROGRAMMERS_TXT = "programmers.txt"
const FOLDER_BOOTLOADERS = "bootloaders"
const FOLDER_CORE = "core"
const FOLDER_CORES = "cores"
const FOLDER_HARDWARE = "hardware"
const FOLDER_LIBRARIES = "libraries"
const FOLDER_PREPROC = "preproc"
const FOLDER_SKETCH = "sketch"
const FOLDER_SYSTEM = "system"
const FOLDER_TOOLS = "tools"
const FOLDER_VARIANTS = "variants"
const hooks_core = hooks + ".core"
const HOOKS_CORE_POSTBUILD = hooks_core + hooks_postbuild_suffix
const HOOKS_CORE_PREBUILD = hooks_core + hooks_prebuild_suffix
const hooks_libraries = hooks + ".libraries"
const HOOKS_LIBRARIES_POSTBUILD = hooks_libraries + hooks_postbuild_suffix
const HOOKS_LIBRARIES_PREBUILD = hooks_libraries + hooks_prebuild_suffix
const hooks_linking = hooks + ".linking"
const HOOKS_LINKING_POSTLINK = hooks_linking + hooks_postlink_suffix
const HOOKS_LINKING_PRELINK = hooks_linking + hooks_prelink_suffix
const hooks_objcopy = hooks + ".objcopy"
const HOOKS_OBJCOPY_POSTOBJCOPY = hooks_objcopy + hooks_postobjcopy_suffix
const HOOKS_OBJCOPY_PREOBJCOPY = hooks_objcopy + hooks_preobjcopy_suffix
const HOOKS_PATTERN_SUFFIX = ".pattern"
const HOOKS_POSTBUILD = hooks + hooks_postbuild_suffix
const hooks_postbuild_suffix = ".postbuild"
const hooks_postlink_suffix = ".postlink"
const hooks_postobjcopy_suffix = ".postobjcopy"
const HOOKS_PREBUILD = hooks + hooks_prebuild_suffix
const hooks_prebuild_suffix = ".prebuild"
const hooks_prelink_suffix = ".prelink"
const hooks_preobjcopy_suffix = ".preobjcopy"
const hooks = "recipe.hooks"
const hooks_sketch = hooks + ".sketch"
const HOOKS_SKETCH_POSTBUILD = hooks_sketch + hooks_postbuild_suffix
const HOOKS_SKETCH_PREBUILD = hooks_sketch + hooks_prebuild_suffix
const IDE_VERSION = "ide_version"
const ID = "_id"
const LIB_CATEGORY_UNCATEGORIZED = "Uncategorized"
const LIB_LICENSE_UNSPECIFIED = "Unspecified"
const LIBRARY_ALL_ARCHS = "*"
const LIBRARY_ARCHITECTURES = "architectures"
const LIBRARY_AUTHOR = "author"
const LIBRARY_CATEGORY = "category"
const LIBRARY_DOT_A_LINKAGE = "dot_a_linkage"
const LIBRARY_EMAIL = "email"
const LIBRARY_FOLDER_ARCH = "arch"
const LIBRARY_FOLDER_SRC = "src"
const LIBRARY_FOLDER_UTILITY = "utility"
const LIBRARY_LICENSE = "license"
const LIBRARY_MAINTAINER = "maintainer"
const LIBRARY_NAME = "name"
const LIBRARY_PARAGRAPH = "paragraph"
const LIBRARY_PROPERTIES = "library.properties"
const LIBRARY_SENTENCE = "sentence"
const LIBRARY_URL = "url"
const LIBRARY_VERSION = "version"
const LOG_LEVEL_DEBUG = "debug"
const LOG_LEVEL_ERROR = "error"
const LOG_LEVEL_INFO = "info"
const LOG_LEVEL_WARN = "warn"
const MSG_ARCH_FOLDER_NOT_SUPPORTED = "'arch' folder is no longer supported! See http://goo.gl/gfFJzU for more information"
const MSG_BOARD_UNKNOWN = "Board {0} (platform {1}, package {2}) is unknown"
const MSG_BOOTLOADER_FILE_MISSING = "Bootloader file specified but missing: {0}"
const MSG_BUILD_OPTIONS_CHANGED = "Build options changed, rebuilding all"
const MSG_CANT_FIND_SKETCH_IN_PATH = "Unable to find {0} in {1}"
const MSG_FQBN_INVALID = "{0} is not a valid fully qualified board name. Required format is targetPackageName:targetPlatformName:targetBoardName."
const MSG_INVALID_QUOTING = "Invalid quoting: no closing [{0}] char found."
const MSG_LIB_LEGACY = "(legacy)"
const MSG_LIBRARIES_MULTIPLE_LIBS_FOUND_FOR = "Multiple libraries were found for \"{0}\""
const MSG_LIBRARIES_NOT_USED = " Not used: {0}"
const MSG_LIBRARIES_USED = " Used: {0}"
const MSG_LIBRARY_CAN_USE_SRC_AND_UTILITY_FOLDERS = "Library can't use both 'src' and 'utility' folders. Double check {0}"
const MSG_LIBRARY_INCOMPATIBLE_ARCH = "WARNING: library {0} claims to run on {1} architecture(s) and may be incompatible with your current board which runs on {2} architecture(s)."
const MSG_LOOKING_FOR_RECIPES = "Looking for recipes like {0}*{1}"
const MSG_MISSING_BUILD_BOARD = "Warning: Board {0}:{1}:{2} doesn''t define a ''build.board'' preference. Auto-set to: {3}"
const MSG_MISSING_CORE_FOR_BOARD = "Selected board depends on '{0}' core (not installed)."
const MSG_MUST_BE_A_FOLDER = "{0} must be a folder"
const MSG_PACKAGE_UNKNOWN = "{0}: Unknown package"
const MSG_PATTERN_MISSING = "{0} pattern is missing"
const MSG_PLATFORM_UNKNOWN = "Platform {0} (package {1}) is unknown"
const MSG_PROGRESS = "Progress {0}"
const MSG_PROP_IN_LIBRARY = "Missing '{0}' from library in {1}"
const MSG_RUNNING_COMMAND = "Ts: {0} - Running: {1}"
const MSG_RUNNING_RECIPE = "Running recipe: {0}"
const MSG_SETTING_BUILD_PATH = "Setting build path to {0}"
const MSG_SKETCH_CANT_BE_IN_BUILDPATH = "Sketch cannot be located in build path. Please specify a different build path"
const MSG_SKIPPING_TAG_ALREADY_DEFINED = "Skipping tag {0} because prototype is already defined"
const MSG_SKIPPING_TAG_BECAUSE_HAS_FIELD = "Skipping tag {0} because it has field {0}"
const MSG_SKIPPING_TAG_WITH_REASON = "Skipping tag {0}. Reason: {1}"
const MSG_UNHANDLED_TYPE_IN_CONTEXT = "Unhandled type {0} in context key {1}"
const MSG_UNKNOWN_SKETCH_EXT = "Unknown sketch file extension: {0}"
const MSG_USING_LIBRARY_AT_VERSION = "Using library {0} at version {1} in folder: {2} {3}"
const MSG_USING_LIBRARY = "Using library {0} in folder: {1} {2}"
const MSG_USING_PREVIOUS_COMPILED_FILE = "Using previously compiled file: {0}"
const MSG_WARNING_LIB_INVALID_CATEGORY = "WARNING: Category '{0}' in library {1} is not valid. Setting to '{2}'"
const MSG_WARNING_PLATFORM_MISSING_VALUE = "Warning: platform.txt from core '{0}' misses property '{1}', using default value '{2}'. Consider upgrading this core."
const MSG_WARNING_PLATFORM_OLD_VALUES = "Warning: platform.txt from core '{0}' contains deprecated {1}, automatically converted to {2}. Consider upgrading this core."
const MSG_WARNING_SPURIOUS_FILE_IN_LIB = "WARNING: Spurious {0} folder in '{1}' library"
const MSG_WRONG_PROPERTIES_FILE = "Property line '{0}' in file {1} is invalid"
const MSG_WRONG_PROPERTIES = "Property line '{0}' is invalid"
const PACKAGE_NAME = "name"
const PACKAGE_TOOLS = "tools"
const PLATFORM_ARCHITECTURE = "architecture"
const PLATFORM_NAME = "name"
const PLATFORM_REWRITE_ADDED = "added"
const PLATFORM_REWRITE_NEW = "new"
const PLATFORM_REWRITE_OLD = "old"
const PLATFORM_URL = "url"
const PLATFORM_VERSION = "version"
const PROGRAMMER_NAME = "name"
const RECIPE_AR_PATTERN = "recipe.ar.pattern"
const RECIPE_C_COMBINE_PATTERN = "recipe.c.combine.pattern"
const RECIPE_C_PATTERN = "recipe.c.o.pattern"
const RECIPE_CPP_PATTERN = "recipe.cpp.o.pattern"
const RECIPE_PREPROC_INCLUDES = "recipe.preproc.includes"
const RECIPE_PREPROC_MACROS = "recipe.preproc.macros"
const RECIPE_S_PATTERN = "recipe.S.o.pattern"
const REWRITING_DISABLED = "disabled"
const REWRITING = "rewriting"
const SPACE = " "
const TOOL_NAME = "name"
const TOOL_URL = "url"
const TOOL_VERSION = "version"
