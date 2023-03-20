package main

import (
	"github.com/bthuilot/resy-booker-go/pkg/cmd"
	"github.com/bthuilot/resy-booker-go/pkg/config"
	"github.com/bthuilot/resy-booker-go/pkg/resy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

func main() {
	config.InitLogger()
	logrus.Info("loading config")
	if err := config.Parse(); err != nil {
		logrus.Fatalf("unable to load config: %s", err)
	}

	logrus.Info("executing command reservations")
	reservations, err := parseReservations()
	if err != nil {
		logrus.Fatalf("unable to parse reservations: %s", err)
	}

	cmd.Execute()

	logrus.Info("constructing Resy client")
	client, err := resy.NewClient()
	if err != nil {
		logrus.Fatalf("unable to construct client: %s", err)
	}

	logrus.Info("constructing Resy booker client")
	booker := resy.NewReservationBooker(client)

	var wg sync.WaitGroup
	for _, req := range reservations {
		wg.Add(1)
		go func(r resy.ReservationRequest) {
			defer wg.Done()
			if booked, err := booker.BookAtMidnight(r); err != nil {
				logrus.Errorf("unable to book reservation: %s", err)
			} else if !booked {
				logrus.Info("unable to book reservation")
			} else {
				logrus.Info("successfully booked reservation")
			}
		}(req)
	}
	wg.Wait()
}

func parseReservations() (reservations []resy.ReservationRequest, err error) {
	err = viper.UnmarshalKey("reservations", &reservations)
	return
}
