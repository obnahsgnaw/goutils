package configutil

import (
	"errors"
	"github.com/jessevdk/go-flags"
	"github.com/obnahsgnaw/goutils/errutil"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"
)

type BaseConfig struct {
	raw               *viper.Viper
	evnPrefix         string
	defConfigFileType string
	defConfigFilePath string
	alias             map[string]string
}

func (s *BaseConfig) ParseFlag(c interface{}) error {
	parser := flags.NewParser(c, flags.Default)
	parser.UnknownOptionHandler = func(option string, arg flags.SplitArgument, args []string) ([]string, error) {
		return nil, nil
	}
	if _, err := parser.Parse(); err != nil {
		var flagErr *flags.Error
		ok := errors.As(err, &flagErr)
		if !ok || !errors.Is(flagErr.Type, flags.ErrHelp) {
			err = configError("config.Parse: config flag parse failed", err)
		}
		return err
	}
	return nil
}

func (s *BaseConfig) ParseFile(c interface{}, file string) error {
	if file != "" || len(os.Args) == 1 {
		s.raw = viper.New()
		if file != "" {
			s.raw.SetConfigFile(file)
		} else {
			s.raw.AddConfigPath(s.getDefaultConfigFilePath())
			s.raw.SetConfigType(s.getDefaultConfigFileType())
		}
		if err := s.raw.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				return configError("config.Parse: config file not found", err)
			} else {
				return configError("config.Parse: config file parse failed", err)
			}
		}
		if err := s.initVar(); err != nil {
			return configError("config.Parse: init var failed", err)
		}

		// evn 需要 PREFIX_KEY 如：USER_APPLICATION.ID
		s.raw.SetEnvPrefix(s.getEnvPrefix())
		s.raw.SetConfigType("env")
		s.raw.AutomaticEnv()

		for k, v := range s.getAlias() {
			s.raw.RegisterAlias(k, strings.ReplaceAll(v, "-", ""))
		}
		if err := s.raw.Unmarshal(c); err != nil {
			return configError("config.Parse: config decode failed", err)
		}
	}

	return nil
}

func (s *BaseConfig) Raw() *viper.Viper {
	return s.raw
}

func (s *BaseConfig) SetAlias(alias map[string]string) {
	s.alias = alias
}

func (s *BaseConfig) getAlias() map[string]string {
	if s.alias == nil {
		s.alias = make(map[string]string)
	}
	return s.alias
}

func (s *BaseConfig) SetEnvPrefix(prefix string) {
	s.evnPrefix = strings.ToUpper(prefix)
}

func (s *BaseConfig) getEnvPrefix() string {
	if s.evnPrefix == "" {
		s.evnPrefix = "DEF"
	}
	return s.evnPrefix
}

func (s *BaseConfig) SetDefaultConfigFileType(ft string) {
	s.defConfigFileType = ft
}

func (s *BaseConfig) getDefaultConfigFileType() string {
	if s.defConfigFileType == "" {
		s.defConfigFileType = "yaml"
	}
	return s.defConfigFileType
}

func (s *BaseConfig) SetDefaultConfigFilePath(fp string) {
	s.defConfigFilePath = fp
}

func (s *BaseConfig) getDefaultConfigFilePath() string {
	if s.defConfigFilePath == "" {
		s.defConfigFilePath = "."
	}
	return s.defConfigFilePath
}

func (s *BaseConfig) initVar() error {
	vars := s.raw.GetStringMap("variable")
	if len(vars) > 0 {
		kvs := s.replaceVar(vars, s.raw.AllSettings(), 0)
		return s.raw.MergeConfigMap(kvs) // set 或写导override中，再获取第一层级key时只会获取override中的数据，造成获取不到config中的，从而遗漏数据
	}
	return nil
}

func (s *BaseConfig) replaceVar(vars map[string]interface{}, kvs map[string]interface{}, level int) map[string]interface{} {
	for k, v := range kvs {
		if level == 0 && k == "variable" {
			continue
		}
		if vv, ok := v.(map[string]interface{}); ok {
			kvs[k] = s.replaceVar(vars, vv, level+1)
		} else {
			vv1, ok1 := v.(string)
			if ok1 {
				if matches := s.find(vv1); len(matches) > 0 {
					for _, match := range matches {
						varKey := strings.Trim(match, "<>")
						varKey = strings.ToLower(varKey)
						if vvv, ok3 := vars[varKey]; ok3 {
							vv1 = strings.Replace(vv1, match, vvv.(string), 1)
						}
					}
					kvs[k] = vv1
				}
			}
		}
	}
	return kvs
}

func (s *BaseConfig) find(text string) []string {
	re := regexp.MustCompile(`<[^>]+>`)
	matches := re.FindAllString(text, -1)
	return matches
}

func configError(msg string, err error) error {
	return errutil.NewFromError(err, msg)
}
