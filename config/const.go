package config

const (
	SETTINGS_FOLDER = "/home/eslap-1118/eslap/config/"
	OPERATION_ENV   = SETTINGS_FOLDER + "operation.yaml"
	NETWORK_ENV     = SETTINGS_FOLDER + "network.yaml"
	DATABASE_ENV    = SETTINGS_FOLDER + "database.yaml"
)

const (
	NO_PRODUCT_TEMPLATE = "no_item"
	NO_TEMPLATE         = "no_template"
)

const (
	EventFrameTxTiming          = -40
	HeartbeatTxTiming           = EventFrameTxTiming - 20
	CfpCutoffTiming             = EventFrameTxTiming - 50
	GrmStackedMsgTxTiming       = EventFrameTxTiming - 20
	ImageGenerationCutoffTiming = EventFrameTxTiming - 120
	TagImageTxTiming            = -25
	TagImageBufferSize          = 16 * 1024
	EventFrameInitAllTxPeriod   = 5
	EventFrameAllTxPeriod       = 3
	EventFrameLongSleepTxPeriod = 15
	UpdateRetryCount            = 5
	UtilizeRf2Mbps              = true
	RfIfs                       = 150
	RfHeaderSize                = 10
	Rf1MbpsTimeUsPerImageByte   = uint32(8)
	Rf2MbpsTimeUsPerImageByte   = uint32(4)
	GmpGrmSpiUnitSize           = 4096
	GrmModemSpiUnitSize         = 256
	SidCountPerGroupMax         = 20 // (GrmModemSpiUnitSize - 6 - 8) / 8
	TagImagePacketSize          = GrmModemSpiUnitSize - 16
	TagImageBlockSize           = TagImagePacketSize * 32
	TagOtaPacketSize            = GrmModemSpiUnitSize - 20
	TagStatusDeleteThreshold    = 5 // number of update fail before tag status delete
	//TagStatusCheckThreshold     = 8  // number of update fail before tag status check
	TagStatusRetryThreshold = 2 // number of update fail before moving tag status to retry
	TaskIdCountMax          = 16
	LongSleepMaxUnitMs      = (1800 * 1000) // 30 minutes
	ScannerModem            = 0
	CfpOffsetBaseMs         = uint32(3)
	RptSlotMs               = 4
	CfpMinMs                = 500
	RptMinMs                = 800
	LongSleepToken          = 0x1A
)

const (
	_ = iota
	COLOR_BW
	COLOR_BWF
	COLOR_BWR
	COLOR_BWY
	COLOR_BWRY
	COLOR_RSVD6
	COLOR_RSVD7
	COLOR_RSVD8
	COLOR_RSVD9
	COLOR_7
)
