package xman

// Package i18n provides an Internationalization and Localization middleware for Macaron applications.
import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/i18n"
	"golang.org/x/text/language"
)

var _localeContextKey = "_locale"
var _defaultOptions = I18nOptions{
	Format:      "%s.ini",
	Directory:   "./config/locale/",
	Langs:       []string{LangZhCN},
	DefaultLang: LangZhCN,
	Names:       []string{"简体中文"},
	Parameter:   "lang",
}

var (
	LangZhCN = "zh-CN"
	LangZhTW = "zh-TW"
)

// isFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func isFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// initLocales initializes language type list and Accept-Language header matcher.
func initLocales(opt I18nOptions) language.Matcher {
	tags := make([]language.Tag, len(opt.Langs))
	for i, lang := range opt.Langs {
		tags[i] = language.Raw.Make(lang)
		fname := fmt.Sprintf(opt.Format, lang)
		// Append custom locale file.
		custom := []interface{}{}
		customPath := path.Join(opt.CustomDirectory, fname)
		if isFile(customPath) {
			custom = append(custom, customPath)
		}

		var locale interface{}
		if data, ok := opt.Files[fname]; ok {
			locale = data
		} else {
			locale = path.Join(opt.Directory, fname)
		}

		err := i18n.SetMessageWithDesc(lang, opt.Names[i], locale, custom...)
		if err != nil && err != i18n.ErrLangAlreadyExist {
			panic(fmt.Errorf("fail to set message file(%s): %v", lang, err))
		}
	}
	return language.NewMatcher(tags)
}

// A Locale describles the information of localization.
type Locale struct {
	i18n.Locale
}

// Language returns language current locale represents.
func (l Locale) Language() string {
	return l.Lang
}

// I18nOptions represents a struct for specifying configuration options for the i18n middleware.
type I18nOptions struct {
	// Suburl of path. Default is empty.
	SubURL string `mapstructure:"sub-url" yaml:"sub-url"`
	// Directory to load locale files. Default is "conf/locale"
	Directory string `mapstructure:"directory" yaml:"directory"`
	// File stores actual data of locale files. Used for in-memory purpose.
	Files map[string][]byte `mapstructure:"files" yaml:"files"`
	// Custom directory to overload locale files. Default is "custom/conf/locale"
	CustomDirectory string `mapstructure:"custom-directory" yaml:"custom-directory"`
	// Langauges that will be supported, order is meaningful.
	Langs []string `mapstructure:"langs" yaml:"langs"`
	// Human friendly names corresponding to Langs list.
	Names []string `mapstructure:"names" yaml:"names"`
	// Default language locale, leave empty to remain unset.
	DefaultLang string `mapstructure:"default-lang" yaml:"default-lang"`
	// Locale file naming style. Default is "locale_%s.ini".
	Format string `mapstructure:"format" yaml:"format"`
	// Name of language parameter name in URL. Default is "lang".
	Parameter string `mapstructure:"parameter" yaml:"parameter"`
	// Redirect when user uses get parameter to specify language.
	Redirect bool `mapstructure:"redirect" yaml:"redirect"`
	// Domain used for `lang` cookie. Default is ""
	CookieDomain string `mapstructure:"cookie-domain" yaml:"cookie-domain"`
}

func newI18n(config I18nOptions) Locale {
	config.SubURL = strings.TrimSuffix(config.SubURL, "/")
	if len(config.Langs) == 0 {
		panic("no language is specified")
	} else if len(config.Langs) != len(config.Names) {
		panic("length of langs is not same as length of names")
	}
	i18n.SetDefaultLang(config.DefaultLang)

	if config.Directory == "" {
		config.Directory = "config/locale"
	}

	if config.Format == "" {
		config.Format = "%s.ini"
	}

	if config.DefaultLang == "" {
		config.DefaultLang = LangZhCN
	}

	return Locale{Locale: i18n.Locale{Lang: config.DefaultLang}}
}

type LangType struct {
	Lang, Name string
}

// I18n is a middleware provides localization layer for your application.
func I18n() gin.HandlerFunc {
	//_globalLocale.
	m := initLocales(_defaultOptions)
	return func(ctx *gin.Context) {
		// 1. Check URL arguments.
		lang := ctx.Query(_defaultOptions.Parameter)

		// 2. Get language information from cookies.
		if lang == "" {
			lang, _ = ctx.Cookie(_defaultOptions.Parameter)
		}

		// 3. Get language information from 'Accept-Language'.
		// The first element in the list is chosen to be the default language automatically.
		if lang == "" {
			tags, _, _ := language.ParseAcceptLanguage(ctx.Request.Header.Get("Accept-Language"))
			tag, _, _ := m.Match(tags...)
			lang = tag.String()
		}

		// Check again in case someone modify by purpose.
		if !i18n.IsExist(lang) {
			lang = _defaultOptions.DefaultLang
		}

		curLang := LangType{
			Lang: lang,
		}

		// Save language information in cookies.
		ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/"+strings.TrimPrefix(_defaultOptions.SubURL, "/"),
			_defaultOptions.CookieDomain, false, false)

		// Set language properties.
		ctx.Set(_localeContextKey, Locale{Locale: i18n.Locale{Lang: lang}})
	}
}
