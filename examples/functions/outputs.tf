output "example" {
  description = "The value of a given environment variable."
  value       = provider::utilities::get_env("GOPATH")
}

output "deep_map" {
  description = "The return value of the deep_merge provider-defined function."
  value = provider::utilities::deep_merge(
    tomap({
      a      = "test",
      b      = "b1",
      c      = "c1",
      banana = "banana"
    }),
    tomap({
      a = "testing",
      b = "b2",
      c = "c2",
      d = "d2"
  }))
}
