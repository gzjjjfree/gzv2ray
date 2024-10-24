package platform

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EnvFlag struct {
	Name    string
	AltName string
}

func NewEnvFlag(name string) EnvFlag { //字符组
	return EnvFlag{
		Name:    name,                   // 输入的字符
		AltName: NormalizeEnvName(name), // 处理后的字符
	}
}

func (f EnvFlag) GetValue(defaultValue func() string) string { //EnvFlag 类型的方法 GetValue 参数为返回 string 的函数，方法返回 string
	if v, found := os.LookupEnv(f.Name); found { // os.LookupEnv 检索由键命名的环境变量的值。如果变量存在于环境中，则返回值（可能为空）且布尔值为真。否则返回值将为空且布尔值为假。
		return v // 有 v2ray.location.confdir 的环境变量，返回变量的值
	}
	if len(f.AltName) > 0 {
		if v, found := os.LookupEnv(f.AltName); found {
			return v
		}
	}

	return defaultValue() // defaultValue 默认返回函数，返回空值
}

func (f EnvFlag) GetValueAsInt(defaultValue int) int {
	useDefaultValue := false
	s := f.GetValue(func() string {
		useDefaultValue = true
		return ""
	})
	if useDefaultValue {
		return defaultValue
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(v)
}

func NormalizeEnvName(name string) string { // 处理路径字符 ReplaceAll 替换字符，ToUpper 转大写，TrimSpace 去除前后空格
	return strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(name)), ".", "_")
}

func getExecutableDir() string {
	exec, err := os.Executable() // Executable 返回启动当前进程的可执行文件的路径名。无法保证该路径仍指向正确的可执行文件。
	// 如果使用符号链接启动进程，则根据操作系统的不同，结果可能是符号链接或它指向的路径。如果需要稳定的结果，path/filepath.EvalSymlinks 可能会有所帮助。
	// 除非发生错误，否则可执行文件会返回绝对路径。	主要用例是查找相对于可执行文件的资源。
	if err != nil {
		return ""
	}
	return filepath.Dir(exec) // Dir 返回路径中除最后一个元素之外的所有元素，通常是路径的目录。删除最后一个元素后，Dir 会在路径上调用 [Clean] 并删除尾部斜杠。
	// 如果路径为空，Dir 将返回“。”。如果路径完全由分隔符组成，Dir 将返回单个分隔符。除非是根目录，否则返回的路径不会以分隔符结尾。
}

func getExecutableSubDir(dir string) func() string {
	return func() string {
		return filepath.Join(getExecutableDir(), dir)
	}
}

func GetPluginDirectory() string {
	const name = "v2ray.location.plugin"
	pluginDir := NewEnvFlag(name).GetValue(getExecutableSubDir("plugins"))
	return pluginDir
}

func GetConfigurationPath() string {
	const name = "v2ray.location.config"
	configPath := NewEnvFlag(name).GetValue(getExecutableDir)
	return filepath.Join(configPath, "config.json")
}

// GetConfDirPath reads "v2ray.location.confdir"
func GetConfDirPath() string {
	const name = "v2ray.location.confdir"
	configPath := NewEnvFlag(name).GetValue(func() string { return "" })
	return configPath
}