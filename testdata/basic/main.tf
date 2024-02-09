resource "terraform_data" "one" {
  input            = "one"
  triggers_replace = "one"
}

# resource "terraform_data" "two" {
#   input = "two"
# }

resource "terraform_data" "three" {
  input = "three-modified"
}
