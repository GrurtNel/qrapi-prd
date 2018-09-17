package scan_history

import (
	"gopkg.in/mgo.v2/bson"
	"qrapi-prd/g/x/web"
	"qrapi-prd/x/logger"
	"qrapi-prd/x/mongodb"
	"qrapi-prd/x/validator"
)

var scanHistoryLog = logger.NewLogger("tbl_scan_history")
var scanHistoryTable = mongodb.NewTable("scan_history", "sh")

type ScanHistory struct {
	mongodb.Model `bson:",inline"`
	CustomerID    string `bson:"customer_id" json:"customer_id"`
	ProductID     string `bson:"product_id" json:"product_id"`
	OrderID       string `bson:"order_id" json:"order_id"`
	URL           string `bson:"url" json:"url"`
	NumberOfScan  int    `bson:"number_of_scan" json:"number_of_scan"`
}

func (scanHistory *ScanHistory) Create() error {
	var existScanHistory *ScanHistory
	scanHistoryTable.FindId(scanHistory.ID).One(&existScanHistory)
	if existScanHistory != nil {
		return scanHistoryTable.UpdateId(scanHistory.ID, bson.M{
			"$inc": bson.M{
				"number_of_scan": 1,
			},
		})
	}
	err := validator.Struct(scanHistory)
	if err != nil {
		scanHistoryLog.Error(err)
		return web.WrapBadRequest(err, "")
	}
	return scanHistoryTable.Insert(scanHistory)
}

func GetByID(id string) *ScanHistory {
	var scanHistory *ScanHistory
	scanHistoryTable.FindId(id).One(&scanHistory)
	if scanHistory == nil {
		return &ScanHistory{
			NumberOfScan: 0,
		}
	}
	return scanHistory
}
