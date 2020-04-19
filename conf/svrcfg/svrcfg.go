/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: svrcfg.go
 * @time: 2017/6/29 9:46
 */
package svrcfg


// 单个服务的配置信息
// type SimpleServerCfg struct {
// }

// 单个服务配置信息，如果存在集群，则需要查看集群节点列表，节点的配置优先于服务配置
// type ServerCfg struct {
// 	Id      string
// 	Name    string // 服务名
// 	MaxProc int
// 	Version string
// 	Memo    string
//
// 	// //共享配置（或主配置）
// 	DbCfgMap    map[string]*DatabaseCfg
// 	RedisCfgMap map[string]*RedisCfg
// 	PostCfg     TcpCfg
// 	ApiCfg      HttpCfg
// 	InnerCfg    TcpCfg
// 	AdditionCfg interface{}
// }
//
// func (p *ServerCfg) Load(cfgPath string) error {
// 	if len(cfgPath) == 0 {
// 		cfgPath = filepath.Dir(os.Args[0])
// 	}
// 	cFile := filepath.Join(cfgPath, "config", fmt.Sprintf("config_%s.toml", env.GetEnvName()))
// 	if _, err := toml.DecodeFile(cFile, p); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (p *ServerCfg) LoadToTarget(cfgPath string, tar interface{}) error {
// 	if len(cfgPath) == 0 {
// 		cfgPath = filepath.Dir(os.Args[0])
// 	}
// 	cFile := filepath.Join(cfgPath, "config", fmt.Sprintf("config_%s.toml", env.GetEnvName()))
// 	if _, err := toml.DecodeFile(cFile, tar); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (p *ServerCfg) GetMaxProc() int {
// 	if p.MaxProc <= 0 || p.MaxProc > runtime.NumCPU() {
// 		return runtime.NumCPU()
// 	}
// 	return p.MaxProc
// }
//
// func (p *ServerCfg) GetVersionStr() string {
// 	return p.Version
// }
//
// func (p *ServerCfg) GetNextVersion() string {
// 	if len(p.Version) <= 0 {
// 		return ""
// 	}
// 	arr := strings.Split(p.Version, ".")
// 	length := len(arr)
// 	num := 0
// 	for i, vs := range arr {
// 		if v, err := strconv.Atoi(vs); err == nil {
// 			num += int(math.Pow10(length-i)) * v
// 		}
// 	}
// 	num += 1

	// arr[0] = strconv.Itoa(num / 100)
	// arr[1] = strconv.Itoa(num % 100 / 10)
	// arr[2] = strconv.Itoa(num % 10)
	//
	// return strings.Join(arr, ".")
// }
//
// func (p *ServerCfg) GetDbCfg(name string) *DatabaseCfg {
// 	if v, ok := p.DbCfgMap[name]; ok {
// 		return v
// 	}
// 	return nil
// }
//
// func (p *ServerCfg) GetRedisCfg(name string) *RedisCfg {
// 	if v, ok := p.RedisCfgMap[name]; ok {
// 		return v
// 	}
// 	return nil
// }
//
// func (p *ServerCfg) FromJson(b []byte) error {
// 	return json.Unmarshal(b, &p)
// }
//
// func (p *ServerCfg) ToJson() ([]byte, error) {
// 	return json.Marshal(p)
// }
