package xman

// Package i18n provides an Internationalization and Localization middleware for Macaron applications.
import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/i18n"
	"golang.org/x/text/language"
)

var _localeContextKey = "_locale"
var _defaultLocale Locale
var _defaultOptions I18nOptions

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
	SubURL string `mapstructure:"sub-url" json:"sub-url" yaml:"sub-url"`
	// Directory to load locale files. Default is "conf/locale"
	Directory string `mapstructure:"directory" json:"directory" yaml:"directory"`
	// File stores actual data of locale files. Used for in-memory purpose.
	Files map[string][]byte `mapstructure:"files" json:"files" yaml:"files"`
	// Custom directory to overload locale files. Default is "custom/conf/locale"
	CustomDirectory string `mapstructure:"custom-directory" json:"custom-directory" yaml:"custom-directory"`
	// Langauges that will be supported, order is meaningful.
	Langs []string `mapstructure:"langs" json:"langs" yaml:"langs"`
	// Human friendly names corresponding to Langs list.
	Names []string `mapstructure:"names" json:"names" yaml:"names"`
	// Default language locale, leave empty to remain unset.
	DefaultLang string `mapstructure:"default-lang" json:"default-lang" yaml:"default-lang"`
	// Locale file naming style. Default is "locale_%s.ini".
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	// Name of language parameter name in URL. Default is "lang".
	Parameter string `mapstructure:"parameter" json:"parameter" yaml:"parameter"`
	// Redirect when user uses get parameter to specify language.
	Redirect bool `mapstructure:"redirect" json:"redirect" yaml:"redirect"`
	// Domain used for `lang` cookie. Default is ""
	CookieDomain string `mapstructure:"cookie-domain" json:"cookie-domain" yaml:"cookie-domain"`
}

func initI18n() {
	opt := sysConf().I18n
	opt.SubURL = strings.TrimSuffix(opt.SubURL, "/")
	if len(opt.Langs) == 0 {
		panic("no language is specified")
	} else if len(opt.Langs) != len(opt.Names) {
		panic("length of langs is not same as length of names")
	}
	i18n.SetDefaultLang(opt.DefaultLang)
	opt.Directory = "config/locale"
	opt.Format = "%s.ini"

	_defaultLocale = Locale{Locale: i18n.Locale{Lang: opt.DefaultLang}}
	_defaultOptions = opt
}

type LangType struct {
	Lang, Name string
}

// I18n is a middleware provides localization layer for your application.
func I18n() gin.HandlerFunc {
	m := initLocales(_defaultOptions)
	return func(ctx *gin.Context) {
		isNeedRedir := false
		hasCookie := false

		// 1. Check URL arguments.
		lang := ctx.Query(_defaultOptions.Parameter)

		// 2. Get language information from cookies.
		if len(lang) == 0 {
			lang, _ = ctx.Cookie("lang")
			hasCookie = true
		} else {
			isNeedRedir = true
		}

		// Check again in case someone modify by purpose.
		if !i18n.IsExist(lang) {
			lang = ""
			isNeedRedir = false
			hasCookie = false
		}

		// 3. Get language information from 'Accept-Language'.
		// The first element in the list is chosen to be the default language automatically.
		if len(lang) == 0 {
			tags, _, _ := language.ParseAcceptLanguage(ctx.Request.Header.Get("Accept-Language"))
			tag, _, _ := m.Match(tags...)
			lang = tag.String()
			isNeedRedir = false
		}

		curLang := LangType{
			Lang: lang,
		}

		// Save language information in cookies.
		if !hasCookie {
			ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/"+strings.TrimPrefix(_defaultOptions.SubURL, "/"),
				_defaultOptions.CookieDomain, false, false)
		}

		restLangs := make([]LangType, 0, i18n.Count()-1)
		langs := i18n.ListLangs()
		names := i18n.ListLangDescs()
		for i, v := range langs {
			if lang != v {
				restLangs = append(restLangs, LangType{v, names[i]})
			} else {
				curLang.Name = names[i]
			}
		}

		// Set language properties.
		ctx.Set(_localeContextKey, Locale{Locale: i18n.Locale{Lang: lang}})
		if _defaultOptions.Redirect && isNeedRedir {
			ctx.Redirect(http.StatusMovedPermanently, _defaultOptions.SubURL+path.Clean(ctx.Request.RequestURI[:strings.Index(ctx.Request.RequestURI, "?")]))
		}
	}
}
