package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"edgetx-to-kml/layout"

	fs "github.com/Egyember/functional-slices"
	"github.com/twpayne/go-kml/v3"
)

var (
	ipath       = flag.String("i", "./log.csv", "log file to convert from")
	altadd      = flag.Int("a", 0, "add this to Altitudes")
	ignoreerror = flag.Bool("f", false, "ignore corrupt records")
	thl         = flag.Int("t", 12, "thread limit")
	barosat     = flag.Float64("b", 0.5, "barometer gps alt ratio to use at specified satelite count")
	satcount    = flag.Int("s", 9, "max satelite count") // not implemented right now
)

const (
	AIR Fmode = iota
	HOR
	ANG
	FAIL
	NONE
)

func tomode(s string) Fmode {
	switch s {
	case "AIR":
		return AIR
	case "HOR":
		return HOR
	case "ANG":
		return ANG
	case "FAIL":
		return FAIL
	}
	return NONE
}

type (
	Fmode  int
	Record struct {
		gps struct {
			lat, lon float64
		}
		sats    int
		time    time.Time
		alt     float64
		baro    float64
		arm     bool
		airmode bool
		mode    Fmode
	}
)

func processRecords(recs []Record) (r []Record) {
	r = make([]Record, len(recs))
	lastalt := float64(0.0)
	for k, v := range recs {
		r[k] = v
		if !v.arm {
			lastalt = r[k].alt
			continue
		}
		r[k].alt = r[k].alt**barosat + r[k].baro*(1-*barosat)
		r[k].alt += lastalt
	}
	return
}

func main() {
	flag.Parse()
	rfd, err := os.Open(*ipath)
	if err != nil {
		panic(err)
	}
	ifd := bufio.NewReader(rfd)
	csv := csv.NewReader(ifd)
	header, err := csv.Read()
	if err != nil {
		panic(err)
	}
	layout := layout.GetLayout(header)
	fmt.Printf("%+v\n", layout)
	data, err := csv.ReadAll()
	if err != nil {
		panic(err)
	}
	parsed := fs.ParMap(data, func(d []string) struct {
		Record
		error
	} {
		p, err := parseRecord(d, layout)
		return struct {
			Record
			error
		}{p, err}
	}, *thl)
	fmt.Println("records parsed: ", len(parsed))
	okrecords := fs.Map(fs.Filter(parsed, func(r struct {
		Record
		error
	},
	) bool {
		return r.error == nil
	}), func(r struct {
		Record
		error
	},
	) Record {
		return r.Record
	})
	if !*ignoreerror {
		if len(okrecords) != len(parsed) {
			panic(errors.New("corrupted records"))
		}
	}

	processed := processRecords(okrecords)
	var place []kml.Element
	var points []kml.Element
	for k, v := range parsed {
		points = append(points, kml.Placemark(
			kml.Name(strconv.Itoa(k)),
			kml.Point(
				kml.Coordinates(kml.Coordinate{Lon: v.gps.lon, Lat: v.gps.lat, Alt: v.alt}),
			),
		))
	}
	place = append(place, kml.Folder(points...))
	place = append(place, kml.Placemark(
		kml.Name("line"),
		kml.LineString(
			kml.Extrude(false),
			kml.Tessellate(true),
			kml.AltitudeMode(kml.AltitudeModeAbsolute),
			kml.Coordinates(fs.Map(processed, func(p Record) (r kml.Coordinate) {
				r.Lat = p.gps.lat
				r.Lon = p.gps.lon
				r.Alt = p.alt
				return
			})...),
		),
	))
	k := kml.KML(kml.Document(place...))
	ofd, err := os.Create("out.kml")
	if err != nil {
		panic(err)
	}
	k.Write(ofd)

	fmt.Println("armed: ", len(fs.Filter(processed, func(r Record) bool {
		return r.arm
	})))
}
