package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed locales/*.json
var localesFS embed.FS

// Language represents a supported language code
type Language string

const (
	LangEnglish    Language = "en-US"
	LangPortuguese Language = "pt-BR"

	DefaultLang Language = LangPortuguese
)

// SupportedLanguages contains all supported language codes
var SupportedLanguages = []Language{LangEnglish, LangPortuguese}

// contextKey is the key type used to store language in context
type contextKey string

const LanguageContextKey contextKey = "language"

// translations holds all loaded translations
var (
	translations = make(map[Language]map[string]string)
	mu           sync.RWMutex
	initialized  bool
)

// Init loads all translation files. Should be called at application startup.
func Init() error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	for _, lang := range SupportedLanguages {
		data, err := localesFS.ReadFile(fmt.Sprintf("locales/%s.json", lang))
		if err != nil {
			return fmt.Errorf("failed to load locale %s: %w", lang, err)
		}

		var rawTrans map[string]interface{}
		if err := json.Unmarshal(data, &rawTrans); err != nil {
			return fmt.Errorf("failed to parse locale %s: %w", lang, err)
		}

		// Flatten nested structure into dot-notation keys
		trans := make(map[string]string)
		flattenTranslations("", rawTrans, trans)

		translations[lang] = trans
	}

	initialized = true
	return nil
}

// flattenTranslations recursively flattens nested maps into dot-notation keys
func flattenTranslations(prefix string, data map[string]interface{}, result map[string]string) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			result[fullKey] = v
		case map[string]interface{}:
			flattenTranslations(fullKey, v, result)
		}
	}
}

// IsSupported checks if a language is supported
func IsSupported(lang Language) bool {
	for _, l := range SupportedLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

// ParseAcceptLanguage parses the Accept-Language header and returns the best match
func ParseAcceptLanguage(header string) Language {
	if header == "" {
		return DefaultLang
	}

	// Split by comma for multiple languages
	parts := strings.Split(header, ",")
	for _, part := range parts {
		// Remove quality values (e.g., "en-US;q=0.9" -> "en-US")
		lang := strings.TrimSpace(strings.Split(part, ";")[0])

		// Try exact match first
		if IsSupported(Language(lang)) {
			return Language(lang)
		}

		// Try base language match (e.g., "en" matches "en-US")
		baseLang := strings.Split(lang, "-")[0]
		for _, supported := range SupportedLanguages {
			if strings.HasPrefix(string(supported), baseLang) {
				return supported
			}
		}
	}

	return DefaultLang
}

// SetLanguage adds the language to the context
func SetLanguage(ctx context.Context, lang Language) context.Context {
	return context.WithValue(ctx, LanguageContextKey, lang)
}

// GetLanguage retrieves the language from the context
func GetLanguage(ctx context.Context) Language {
	if lang, ok := ctx.Value(LanguageContextKey).(Language); ok {
		return lang
	}
	return DefaultLang
}

// T translates a key to the language in the context
func T(ctx context.Context, key string) string {
	lang := GetLanguage(ctx)
	return Translate(lang, key)
}

// Translate translates a key to the specified language
func Translate(lang Language, key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if trans, ok := translations[lang]; ok {
		if msg, ok := trans[key]; ok {
			return msg
		}
	}

	// Fallback to default language
	if trans, ok := translations[DefaultLang]; ok {
		if msg, ok := trans[key]; ok {
			return msg
		}
	}

	// Return key if no translation found
	return key
}

// TWithParams translates a key and replaces placeholders with values
// Placeholders are in the format {key}
func TWithParams(ctx context.Context, key string, params map[string]string) string {
	msg := T(ctx, key)
	for k, v := range params {
		msg = strings.ReplaceAll(msg, "{"+k+"}", v)
	}
	return msg
}

// TranslateWithParams translates with the specified language and replaces placeholders
func TranslateWithParams(lang Language, key string, params map[string]string) string {
	msg := Translate(lang, key)
	for k, v := range params {
		msg = strings.ReplaceAll(msg, "{"+k+"}", v)
	}
	return msg
}

// GetAllTranslations returns all translations for a language (useful for bulk operations)
func GetAllTranslations(lang Language) map[string]string {
	mu.RLock()
	defer mu.RUnlock()

	if trans, ok := translations[lang]; ok {
		// Return a copy to prevent modifications
		result := make(map[string]string, len(trans))
		for k, v := range trans {
			result[k] = v
		}
		return result
	}
	return nil
}

// GetTranslationsByPrefix returns all translations with keys starting with the prefix
func GetTranslationsByPrefix(lang Language, prefix string) map[string]string {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string]string)
	if trans, ok := translations[lang]; ok {
		for k, v := range trans {
			if strings.HasPrefix(k, prefix) {
				result[k] = v
			}
		}
	}
	return result
}
