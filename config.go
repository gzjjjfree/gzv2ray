//go:build !confonly
// +build !confonly

package core

import (
	"io"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/gzjjjfree/gzv2ray/common"
	"github.com/gzjjjfree/gzv2ray/common/buf"
	"github.com/gzjjjfree/gzv2ray/common/cmdarg"
	"github.com/gzjjjfree/gzv2ray/main/confloader"
)

// ConfigFormat is a configurable format of V2Ray config file.
type ConfigFormat struct { // ConfigFormat 是 V2Ray 配置文件的可配置格式。
	Name      string
	Extension []string
	Loader    ConfigLoader
}

// ConfigLoader is a utility to load V2Ray config from external source.
type ConfigLoader func(input interface{}) (*Config, error) // ConfigLoader 是一个从外部源加载 V2Ray 配置的实用程序。

var (
	configLoaderByName = make(map[string]*ConfigFormat)
	configLoaderByExt  = make(map[string]*ConfigFormat)
)

// RegisterConfigLoader add a new ConfigLoader. 添加一个新的 ConfigLoader
func RegisterConfigLoader(format *ConfigFormat) error {
	name := strings.ToLower(format.Name)
	if _, found := configLoaderByName[name]; found {
		return newError(format.Name, " already registered.")
	}
	configLoaderByName[name] = format

	for _, ext := range format.Extension {
		lext := strings.ToLower(ext)
		if f, found := configLoaderByExt[lext]; found {
			return newError(ext, " already registered to ", f.Name)
		}
		configLoaderByExt[lext] = format
	}

	return nil
}

func getExtension(filename string) string {
	idx := strings.LastIndexByte(filename, '.') // LastIndexByte 返回 s 中 c 的最后一个实例的索引，如果 c 不存在于 s 中，则返回 -1，这里读取路径中的文件名后缀前的符号 .
	if idx == -1 { // 如果没有，返回空
		return ""
	}
	return filename[idx+1:]   // . 后面是后缀名
}

// LoadConfig loads config with given format from given source. 
// LoadConfig 从给定源加载具有给定格式的配置
// input accepts 2 different types: 
// 输入接受两种不同的类型：
// * []string slice of multiple filename/url(s) to open to read
// * []string 切片，包含多个文件名/url，用于打开并读取
// * io.Reader that reads a config content (the original way)
// io.Reader 读取配置内容（原始方式）
func LoadConfig(formatName string, filename string, input interface{}) (*Config, error) { // formatName 文件的类型，filename 文件名切片，input 通过 io.Reader 读取
	ext := getExtension(filename)  // ext 找 filename 后缀名
	if len(ext) > 0 {
		if f, found := configLoaderByExt[ext]; found { // 如果是 v2ray 的可配置格式
			return f.Loader(input)  // 返回配置加载函数，把文件加载到接口中
		}
	}

	if f, found := configLoaderByName[formatName]; found {   // 通过 formatName 类型确认格式
		return f.Loader(input)
	}

	return nil, newError("Unable to load config in ", formatName).AtWarning()
}

func loadProtobufConfig(data []byte) (*Config, error) {
	config := new(Config)
	if err := proto.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func init() {
	common.Must(RegisterConfigLoader(&ConfigFormat{
		Name:      "Protobuf",
		Extension: []string{"pb"},
		Loader: func(input interface{}) (*Config, error) {
			switch v := input.(type) {
			case cmdarg.Arg:
				r, err := confloader.LoadConfig(v[0])
				common.Must(err)
				data, err := buf.ReadAllToBytes(r)
				common.Must(err)
				return loadProtobufConfig(data)
			case io.Reader:
				data, err := buf.ReadAllToBytes(v)
				common.Must(err)
				return loadProtobufConfig(data)
			default:
				return nil, newError("unknow type")
			}
		},
	}))
}


