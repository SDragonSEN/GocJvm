package startup

import "flag"

type Cmd struct {
	Version   bool
	Help      bool
	Classpath string
	MainClass string
	Args      []string
	Jar       string
}

var CmdPara Cmd

func ParseCmd() {
	flag.BoolVar(&CmdPara.Help, "help", false, "帮助信息")
	flag.BoolVar(&CmdPara.Help, "?", false, "帮助信息")
	flag.BoolVar(&CmdPara.Version, "version", false, "版本信息")
	flag.StringVar(&CmdPara.Classpath, "cp", "", "类搜索路径")
	flag.StringVar(&CmdPara.Classpath, "classpath", "", "类搜索路径")
	flag.StringVar(&CmdPara.Jar, "jar", "", "jar加载路径")

	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		if CmdPara.Jar != "" {
			//参数
			CmdPara.Args = args[:]
		} else {
			//主类路径
			CmdPara.MainClass = args[0]
			//参数
			CmdPara.Args = args[1:]
		}
	}
}
