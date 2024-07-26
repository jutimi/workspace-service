package utils

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/gosimple/unidecode"
)

func ConvertToUpperCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "_")
	return strings.ToUpper(str)
}

func ConvertToCamelCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, " ")
	parts := strings.Split(str, " ")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	// Concatenate the words
	result := strings.Join(parts, "")
	return result
}

func ConvertToSnakeCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "_")
	return strings.ToLower(str)
}

func Slugify(str string) string {
	decodeStr := unidecode.Unidecode(str)
	return slug.MakeLang(decodeStr, "vi")
}

func ConvertStringToUUID(str string) (uuid.UUID, error) {
	uuid, err := uuid.Parse(str)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func ConvertSliceStringToUUID(str []string) ([]uuid.UUID, error) {
	var result []uuid.UUID
	for _, data := range str {
		convertData, err := ConvertStringToUUID(data)
		if err != nil {
			return nil, err
		}
		result = append(result, convertData)
	}
	return result, nil
}
