package config

type OperationEnv struct {
	ConfigDir         string `yaml:"configDir"`
	CertFile          string `yaml:"certFile"`
	LogDir            string `yaml:"logDir"`
	TmpDir            string `yaml:"tmpDir"`
	TagImageDir       string `yaml:"tagImageDir"`
	TemplateDir       string `yaml:"templateDir"`
	FontDir           string `yaml:"fontDir"`
	ImageDir          string `yaml:"imageDir"`
	ExternalBinaryDir string `yaml:"externalBinaryDir"`
	DataBackupDir     string `yaml:"dataBackupDir"`

	LogLevel                  string `yaml:"logLevel"`
	ApBrokerRetryTimingSecond int    `yaml:"apBrokerRetryTimingSecond"`
	ImgGenThreadCount         int    `yaml:"imageGenerationThreadCount"`
	ImgGenReqPort             string `yaml:"imageGenerationReqPort"`
	ImgGenRespPort            string `yaml:"imageGenerationRespPort"`
	ImgGenPubPort             string `yaml:"imageGenerationPubPort"`
	EslApTimerReqPort         string `yaml:"eslapTimerReqPort"`
	DeassignThreadCount       int    `yaml:"deassignThreadCount"`
	FontFacePreloadCount      int    `yaml:"fontFacePreloadCount"`

	LastTaskIdBackupFile  string `yaml:"lastTaskIdBackupFile"`
	ProductDataBackupFile string `yaml:"productDataBackupFile"`
	AssignDataBackupFile  string `yaml:"assignDataBackupFile"`
	NfcDataBackupFile     string `yaml:"nfcDataBackupFile"`
	EventFrameTxTiming    int    `yaml:"eventFrameTxTiming"`
	TagImageTxTimging     int    `yaml:"tagImageTxTimging"`
	ScanProfile           [6]int `yaml:"scanProfile"`
	BackoffBase           int    `yaml:"backoffBase"`
	BackoffMulFactor      int    `yaml:"backoffMulfactor"`
	FreezerTagMultiplier  int    `yaml:"freezerTagMultiplier"`
	TagDistributionMinute int    `yaml:"tagDistributionMinute"`
	PageRotationMacPage   int    `yaml:"pageRotationMacPage"`

	LogMaxSizeMb  int  `yaml:"logMaxSizeMb"`
	LogMaxBackup  int  `yaml:"logMaxBackup"`
	LogMaxAgeDays int  `yaml:"logMaxAgeDays"`
	LogCompress   bool `yaml:"logCompress"`
	GRpcMaxSize   int  `yaml:"gRpcMaxSize"`
}

type DatabaseEnv struct {
	Driver            string `yaml:"driver"`
	FileName          string `yaml:"fileName"`
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	User              string `yaml:"user"`
	Passwd            string `yaml:"passwd"`
	Ssl               string `yaml:"ssl"`
	IpcKey            string `yaml:"ipcKey"`
	Log               string `yaml:"log"`
	DatabaseDir       string `yaml:"databaseDir"`
	DatabaseBackupDir string `yaml:"databaseBackupDir"`
}

type NetworkEnv struct {
	SiteId            string   `yaml:"siteId"`
	SiteCode          uint32   `yaml:"siteCode"`
	StoreCode         string   `yaml:"storeCode"`
	Ip                string   `yaml:"ip"`
	DefaultGwIP       string   `yaml:"defaultGwIp"`
	Netmask           string   `yaml:"netmask"`
	NameServers       []string `yaml:"nameServers"`
	TimeZone          string   `yaml:"timeZone"`
	TimeServerUrls    []string `yaml:"timeServerUrls"`
	InterApPort       string   `yaml:"interApPort"`
	InterApPortTarget string   `yaml:"interApPortTarget"`
	ApBrokerUrl       string   `yaml:"apBrokerUrl"`
	EthernetInterface string   `yaml:"ethernetInterface"`
}

type GrmEnv struct {
	SpiDevice  string `yaml:"spiDevice"`
	SpiMode    string `yaml:"spiMode"`
	SpiSpeedHz uint32 `yaml:"spiSpeedHz"`
	SpiPinCs   uint32 `yaml:"spiPinCs"`
	SpiPinMosi uint32 `yaml:"spiPinMosi"`
	SpiPinMiso uint32 `yaml:"spiPinMiso"`
	SpiPinClk  uint32 `yaml:"spiPinClk"`
}

type LastTaskIdBackup struct {
	LastTaskId uint32 `yaml:"lastTaskId"`
}
