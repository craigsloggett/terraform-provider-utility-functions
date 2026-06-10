output "example" {
  description = "The result of deeply merging two maps."
  value = provider::utilities::deep_merge(
    tomap({ a = "test", b = "b1", c = "c1" }),
    tomap({ a = "testing", b = "b2", c = "c2", d = "d2" })
  )
}
