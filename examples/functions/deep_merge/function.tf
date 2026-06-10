output "example" {
  description = "The result of deeply merging two nested objects."
  value = provider::utilities::deep_merge(
    {
      ports = { http = 80, https = 443 }
      tags  = { environment = "development", team = "platform" }
    },
    {
      tags = { environment = "production" }
    }
  )
}
