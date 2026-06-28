package main

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"edgetx-to-kml/layout"
)

func parseRecord(row []string, l layout.Layout) (r Record, err error) {
	err = nil
	r.time, err = time.Parse("2006-01-02 15:04:05.000", row[l.Date]+" "+row[l.Time])
	if err != nil {
		return
	}
	parts := strings.Split(row[l.Gps], " ")
	if len(parts) != 2 {
		err = errors.New("gps error:" + row[l.Gps] + "parts:" + strconv.Itoa(len(parts)))
		return
	}
	r.gps.lat, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return
	}
	r.gps.lon, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return
	}

	r.alt, err = strconv.ParseFloat(row[l.Alt], 64)
	if err != nil {
		return
	}
	r.alt += float64(*altadd)

	r.baro, err = strconv.ParseFloat(row[l.Baro], 64)
	if err != nil {
		return
	}

	satint, err := strconv.ParseInt(row[l.Sats], 10, 64)
	if err != nil {
		return
	}

	r.sats = int(satint)

	if row[l.Mode][len(row[l.Mode])-1] != '*' {
		r.arm = true
		r.mode = tomode(row[l.Mode])
	} else {
		r.arm = false
		r.mode = tomode(row[l.Mode][:len(row[l.Mode])-1])
	}

	// todo: fix this
	r.airmode = true

	return
}
