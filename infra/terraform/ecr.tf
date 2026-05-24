variable "ecr_repos" {
  type = set(string)
  default = [
    "backend/auth",
    "backend/chat",
    "backend/expo-notification",
    "backend/fcmnotification",
    "backend/messagepersist",
    "backend/messagepreprocess",
    "backend/messagerelay",
    "backend/notification",
    "backend/onlineconversation",
    "backend/signalrelay"
  ]
}

resource "aws_ecr_repository" "repos" {
  for_each             = var.ecr_repos
  name                 = each.value
  image_tag_mutability = "MUTABLE"
}

resource "aws_ecr_lifecycle_policy" "repo_cleanup" {
  for_each   = aws_ecr_repository.repos
  repository = each.value.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 3 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 3
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}