variable "map1" {
  type        = map(string)
  description = "A map to use for the input to the deep_merge function."
  default = {
    a = "test"
    b = "b1"
    c = "c1"
  }
}

variable "map2" {
  type        = map(string)
  description = "A map to use for the input to the deep_merge function."
  default = {
    a = "testing"
    b = "b2"
    c = "c3"
  }
}

variable "object1" {
  type        = object({})
  description = "An object to use for the input to the deep_merge function."
  default = {
    a = "test"
    b = {
      b1 = "b1"
      b2 = "b2"
      b4 = {
        b5 = "b5"
      }
    }
    c = [1, 2, 3]
  }
}

variable "object2" {
  type        = object({})
  description = "An object to use for the input to the deep_merge function."
  default = {
    b = {
      b1 = "b1-updated"
      b3 = "b3"
      b4 = {
        b5 = "b5-updated"
      }
    }
    c = [1, 2, 3, 4]
  }
}
