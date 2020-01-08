package option

import "github.com/spf13/pflag"

type Config struct {
	DBConnectionInfo  string
	DBType            string
}

type Spider struct {
	Config
} 

func NewSpider() *Spider {
	return &Spider{}
}

func (s *Spider) AddFlags(fs *pflag.FlagSet)  {
	fs.StringVar(&s.Config.DBConnectionInfo, "db", "root:123456@tcp(127.0.0.1:3306)/stock", "DB Connection Info")
	fs.StringVar(&s.Config.DBType, "dbType", "mysql", "DB Type")
}
