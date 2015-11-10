package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"strconv"
	"strings"
)

type CTagsToPrototypes struct{}

func (s *CTagsToPrototypes) Run(context map[string]interface{}) error {
	tags := context[constants.CTX_COLLECTED_CTAGS].([]map[string]string)

	lineWhereToInsertPrototypes, err := findLineWhereToInsertPrototypes(tags)
	if err != nil {
		return utils.WrapError(err)
	}
	if lineWhereToInsertPrototypes != -1 {
		context[constants.CTX_LINE_WHERE_TO_INSERT_PROTOTYPES] = lineWhereToInsertPrototypes
	}

	prototypes := toPrototypes(tags)
	context[constants.CTX_PROTOTYPES] = prototypes

	return nil
}

func findLineWhereToInsertPrototypes(tags []map[string]string) (int, error) {
	firstFunctionLine, err := firstFunctionAtLine(tags)
	if err != nil {
		return -1, utils.WrapError(err)
	}
	firstFunctionPointerAsArgument, err := firstFunctionPointerUsedAsArgument(tags)
	if err != nil {
		return -1, utils.WrapError(err)
	}
	if firstFunctionLine != -1 && firstFunctionPointerAsArgument != -1 {
		if firstFunctionLine < firstFunctionPointerAsArgument {
			return firstFunctionLine, nil
		} else {
			return firstFunctionPointerAsArgument, nil
		}
	} else if firstFunctionLine == -1 {
		return firstFunctionPointerAsArgument, nil
	} else {
		return firstFunctionLine, nil
	}
}

func firstFunctionPointerUsedAsArgument(tags []map[string]string) (int, error) {
	functionNames := collectFunctionNames(tags)
	for _, tag := range tags {
		if functionNameUsedAsFunctionPointerIn(tag, functionNames) {
			return strconv.Atoi(tag[FIELD_LINE])
		}
	}
	return -1, nil
}

func functionNameUsedAsFunctionPointerIn(tag map[string]string, functionNames []string) bool {
	for _, functionName := range functionNames {
		if strings.Index(tag[FIELD_CODE], "&"+functionName) != -1 {
			return true
		}
	}
	return false
}

func collectFunctionNames(tags []map[string]string) []string {
	names := []string{}
	for _, tag := range tags {
		if tag[FIELD_KIND] == KIND_FUNCTION {
			names = append(names, tag[constants.CTAGS_FIELD_FUNCTION_NAME])
		}
	}
	return names
}

func firstFunctionAtLine(tags []map[string]string) (int, error) {
	for _, tag := range tags {
		if !tagIsUnknown(tag) && !tagHasAtLeastOneField(tag, FIELDS_MARKING_UNHANDLED_TAGS) && tag[FIELD_KIND] == KIND_FUNCTION {
			return strconv.Atoi(tag[FIELD_LINE])
		}
	}
	return -1, nil
}

func toPrototypes(tags []map[string]string) []*types.Prototype {
	prototypes := []*types.Prototype{}
	for _, tag := range tags {
		if tag[FIELD_SKIP] != TRUE {
			ctag := types.Prototype{FunctionName: tag[constants.CTAGS_FIELD_FUNCTION_NAME], Prototype: tag[KIND_PROTOTYPE], Modifiers: tag[KIND_PROTOTYPE_MODIFIERS], Line: tag[FIELD_LINE], Fields: tag}
			prototypes = append(prototypes, &ctag)
		}
	}
	return prototypes
}
