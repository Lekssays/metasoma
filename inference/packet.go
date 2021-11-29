package main

import (
	"encoding/csv"
    "log"
    "os"
	"strconv"
)

type Packet struct {
	timestamp int64
	srcIp     int32
	dstIp     int32
	features  []float32
}

func ReadCSV(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}


func ParsePackets(records [][]string) []Packet {
	var packets []Packet

	featuresSize := len(records[0])
	for i := 1; i < len(records); i++ {
		var tmpPacket Packet
		tmpPacket.timestamp, _ = strconv.ParseInt(records[i][0], 10, 64)
		var srcIp int64
		var dstIp int64
		srcIp, _ = strconv.ParseInt(records[i][1], 10, 32)
		dstIp, _ = strconv.ParseInt(records[i][2], 10, 32)
		tmpPacket.srcIp = int32(srcIp)
		tmpPacket.dstIp = int32(dstIp)
		for j := 3; j < featuresSize; j++ {
			var tmpFeature float64
			tmpFeature, _ = strconv.ParseFloat(records[i][j], 32)
			tmpPacket.features = append(tmpPacket.features, float32(tmpFeature))
		}
		packets = append(packets, tmpPacket)
	}

	return packets
}
