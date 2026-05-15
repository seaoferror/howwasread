# variable "ecr_repos" {
#   type = set(string)
#   default = [
#     "app-service-1",
#     "app-service-2",
#     "background-worker",
#     # TODO: add all 10 names here
#   ]
# }
#
# resource "aws_ecr_repository" "repos" {
#   for_each             = var.ecr_repos
#   name                 = each.value
#   image_tag_mutability = "IMMUTABLE"
# }
#
# resource "aws_ecr_lifecycle_policy" "repo_cleanup" {
#   for_each   = aws_ecr_repository.repos
#   repository = each.value.name
#
#   policy = jsonencode({
#     rules = [
#       {
#         rulePriority = 1
#         description  = "Keep last 30 images"
#         selection = {
#           tagStatus   = "any"
#           countType   = "imageCountMoreThan"
#           countNumber = 3
#         }
#         action = {
#           type = "expire"
#         }
#       }
#     ]
#   })
# }