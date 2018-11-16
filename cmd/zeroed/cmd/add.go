package cmd

import (
	"fmt"
	"os"

	"github.com/bitnami-labs/zerozone/pkg/model"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	recordTTL uint32
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new DNS record",
	Long:  `Add a new DNS record`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, save, err := openZone(viper.GetString(fileCfg))
		if err != nil {
			return err
		}
		r, ok := zone.FindRecord(recordName, recordType)
		if !ok {
			zone.Records = append(zone.Records, model.ResourceRecordSet{
				Name: recordName,
				Type: recordType,
			})
			r = &zone.Records[len(zone.Records)-1]
		}
		r.TTL = recordTTL

		found := false
		for _, d := range r.RRDatas {
			if d == recordData {
				found = true
			}
		}
		if !found {
			r.RRDatas = append(r.RRDatas, recordData)
		}

		zoneName, err := getZoneName()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "\n")
		fqdn := fmt.Sprintf("%s.%s.%s", recordName, zoneName, viper.GetString(zeroZoneDomainCfg))
		fmt.Printf("%s %d IN %s %s\n", fqdn, recordTTL, recordType, recordData)
		fmt.Fprintf(os.Stderr, "\n")

		if err := save(); err != nil {
			return err
		}

		return waitForRecord(fqdn, recordType, zoneName)
	},
}

func waitForRecord(fqdn string, recordType string, zoneName string) error {
	t, ok := dns.StringToType[recordType]
	if !ok {
		return fmt.Errorf("unknown record type %q", recordType)
	}

	fmt.Println("waiting until DNS server resolves our new record")

	for {
		m1 := new(dns.Msg)
		m1.Id = dns.Id()
		m1.RecursionDesired = true
		m1.Question = make([]dns.Question, 1)
		m1.Question[0] = dns.Question{fmt.Sprintf("%s.", fqdn), t, dns.ClassINET}

		c := new(dns.Client)
		in, _, err := c.Exchange(m1, "0zone.ns.mkm.pub:53")
		if err == nil {
			fmt.Println("record resolved")
			fmt.Printf("%s\n", in)
			return nil
		}
	}
}

func init() {
	rootCmd.AddCommand(addCmd)

	registerEditFlags(addCmd)

	addCmd.Flags().Uint32Var(&recordTTL, "ttl", 60, "record ttl")
	addCmd.MarkFlagRequired("data")
}
