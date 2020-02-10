package provider

import (
	"bytes"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func splitLen(s string, n int) []string {
	sub := ""
	subs := []string{}
	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}

func keys(data map[string]string) []string {
	var result = []string{}
	for k := range data {
		result = append(result, k)
	}
	return result
}

func mergeSchemas(a map[string]*schema.Schema, b map[string]*schema.Schema) map[string]*schema.Schema {
	result := a
	for k, v := range b {
		result[k] = v
	}
	return result
}

// Converts duration string to an int representing the number of seconds, which is used by the Keycloak API
// Ex: "1h" => 3600
func getSecondsFromDurationString(s string) (int, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}

	return int(duration.Seconds()), nil
}

// Converts number of seconds from Keycloak API to a duration string used by the provider
// Ex: 3600 => "1h0m0s"
func getDurationStringFromSeconds(seconds int) string {
	return (time.Duration(seconds) * time.Second).String()
}

// This will suppress the Terraform diff when comparing duration strings.
// As long as both strings represent the same number of seconds, it makes no difference to the Keycloak API
func suppressDurationStringDiff(_, old, new string, _ *schema.ResourceData) bool {
	if old == "" || new == "" {
		return false
	}

	oldDuration, _ := time.ParseDuration(old)
	newDuration, _ := time.ParseDuration(new)

	return oldDuration.Seconds() == newDuration.Seconds()
}

func handleNotFoundError(err error, data *schema.ResourceData) error {
	if keycloak.ErrorIs404(err) {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer exists", data.Id())
		data.SetId("")

		return nil
	}

	return err
}

func interfaceSliceToStringSlice(iv []interface{}) []string {
	var sv []string
	for _, i := range iv {
		sv = append(sv, i.(string))
	}

	return sv
}
