package main

import (
	"context"
	"log"

	"github.com/craigsloggett/terraform-provider-utility-functions/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/craigsloggett/utility-functions",
		Debug:   false,
	}

	err := providerserver.Serve(context.Background(), provider.NewUtilityFunctionsProvider(), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
