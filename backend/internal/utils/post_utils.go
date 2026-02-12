package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
)

func GenerateSlug(title string)string{
	slug:=strings.ToLower(title)
	reg, _ := regexp.Compile("[^a-z0-9 ]+")
	slug = reg.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	return strings.Trim(slug, "-")
}

func GenerateUniqueSlug(title string) string {
	base := GenerateSlug(title)
	b := make([]byte, 2)
	rand.Read(b)
	return fmt.Sprintf("%s-%x", base, b)
}