output "example" {
  description = "The value of a given environment variable."
  value       = provider::utilities::get_env("GOPATH")
}

output "deep_map" {
  description = "The return value of the deep_merge provider-defined function."
  value       = provider::utilities::deep_merge(var.map1, var.map2)
}

# Expected Output:

# deep_map = {
#   "a" = "test"
#   "b" = {
#     "b1" = "b1-updated"
#     "b2" = "b2"
#     "b3" = "b3"
#     "b4" = {
#       "b5" = "b5-updated"
#     }
#   }
#   "c" = [
#     1,
#     2,
#     3,
#     4,
#   ]
# }
