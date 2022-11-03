package main

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func main() {
	// https://github.com/sapics/ip-location-db/
	resp, err := http.Get(
		"https://raw.githubusercontent.com/sapics/ip-location-db/master/geolite2-city/geolite2-city-ipv4-num.csv.gz",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("could not pull data from ip-location-db, status code is %d", resp.StatusCode)
	}

	f, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer client.Close()

	pipe := client.Pipeline()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	var last uint32 = 0
	var i uint32 = 0
	for _, line := range lines {
		start, _ := strconv.ParseUint(line[0], 10, 32)
		end, _ := strconv.ParseUint(line[1], 10, 32)

		if uint32(start) > last+1 {
			fakeCity := fmt.Sprintf("%d|%s|%s", uint32(start)-1, "", "")
			err = pipe.ZAdd(
				ctx,
				"ipv4",
				&redis.Z{Score: float64(start - 1), Member: fakeCity},
			).Err()
			if err != nil {
				log.Fatal(err)
			}
			i++
		}
		last = uint32(end)

		realCity := fmt.Sprintf("%d|%s|%s", uint32(end), line[2], line[5])
		err = pipe.ZAdd(
			ctx,
			"ipv4",
			&redis.Z{Score: float64(end), Member: realCity},
		).Err()

		if err != nil {
			log.Fatal(err)
		}
		i++

		if i%10000 == 0 {
			_, err = pipe.Exec(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
