package cmd

import (
	"fmt"
	"github.com/bthuilot/book-it/pkg/booker"
	"github.com/bthuilot/book-it/pkg/config"
	"github.com/bthuilot/book-it/pkg/resy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

var args struct {
	date       *string
	venueID    *int
	partySize  *int
	types      *[]string
	time       *string
	timeSpread *int
}

func init() {
	args.venueID = rootCmd.PersistentFlags().IntP("venue-id", "v", 0, "ID of the venue you want to reserve from ")

	args.date = rootCmd.PersistentFlags().StringP("date", "d", "", "The date of when to book in 'mm/dd/yy' format")

	args.partySize = rootCmd.PersistentFlags().IntP("party-size", "p", 0, "the amount of people the reservation should be for")

	args.types = rootCmd.PersistentFlags().StringArray("include-types", nil, "filter for what types of reservations to include, i.e. Booth, Outdoor, etc. Must match exact name shown on Resy")
	args.time = rootCmd.PersistentFlags().StringP("time", "t", "", "The time to make the reservation for in kitchen time format (i.e. 3:04PM)")

	args.timeSpread = rootCmd.PersistentFlags().IntP("time-spread", "s", 0, "The 'spread' of acceptable times. When set, the program will consider times within this duration of seconds from the targeted time to be acceptable reservations")

}

func parseArgs() (venueID int, partySize int, date time.Time, spread time.Duration, includeTypes []string, err error) {
	if args.venueID != nil && *args.venueID != -1 {
		venueID = *args.venueID
	}
	if args.date != nil && args.time != nil {
		date, err = time.Parse("01/02/06 3:04PM", fmt.Sprintf("%s %s", *args.date, *args.time))
	}
	if args.partySize != nil && *args.partySize != -1 {
		partySize = *args.partySize
	}

	if args.timeSpread != nil {
		spread = time.Duration(*args.timeSpread) * time.Second
	}

	if args.types != nil {
		includeTypes = *args.types
	}
	return
}

var rootCmd = &cobra.Command{
	Use:   "book-it",
	Short: "book-it will book restaurant reservations using the Resy API",
	Long:  `A CLI tool that can be used to book hard-to-get restaurant reservations using the resy API`,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		config.InitLogger()
		config.ParseEnv()
		if err = cmd.MarkFlagRequired("venue-id"); err != nil {
			return fmt.Errorf("venue ID must be set")
		}
		if err = cmd.MarkFlagRequired("date"); err != nil {
			return fmt.Errorf("date must be set")
		}

		if err = cmd.MarkFlagRequired("party-size"); err != nil {
			err = fmt.Errorf("party size must be set")
			return
		}
		if err = cmd.MarkFlagRequired("time"); err != nil {
			err = fmt.Errorf("time must be set")
			return
		}
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		venueID, partySize, date, spread, includeTypes, err := parseArgs()
		if err != nil {
			logrus.Fatalf("unable to parse args: %s", err)
		}
		logrus.Info("constructing Resy client")
		client, err := resy.NewClient(
			resy.WithCredentialsOpts(viper.GetString("email"), viper.GetString("password")),
		)
		if err != nil {
			logrus.Fatalf("unable to construct client: %s", err)
		}

		logrus.Info("constructing Resy booker client")
		booker := booker.NewBooker(client)
		reservation, err := booker.Book(venueID, partySize, date, spread, includeTypes...)
		if err != nil {
			logrus.Fatalf("unable to book reservation for venue %d on %s: %s", venueID, date, err)
		}
		fmt.Printf("reservation: %+v\n", reservation)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
