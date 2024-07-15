# provider
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  profile = var.profile
  region  = var.region
}

# import
data "terraform_remote_state" "stinkyfingers" {
  backend = "s3"
  config = {
    bucket  = "remotebackend"
    key     = "stinkyfingers/terraform.tfstate"
    region  = "us-west-1"
    profile = var.profile
  }
}

# Lambda
resource "aws_lambda_permission" "server" {
  statement_id  = "AllowExecutionFromApplicationLoadBalancer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.server.arn
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.target.arn
}

resource "aws_lambda_permission" "server_live" {
  statement_id  = "AllowExecutionFromApplicationLoadBalancer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_alias.server_live.arn
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.target.arn
}

resource "aws_lambda_alias" "server_live" {
  name             = "live"
  description      = "set a live alias"
  function_name    = aws_lambda_function.server.arn
  function_version = aws_lambda_function.server.version
}

resource "aws_lambda_function" "server" {
  filename         = "../lambda.zip"
  function_name    = "shoppinglistapi"
  role             = aws_iam_role.lambda_role.arn
  handler          = "lambda-lambda"
  runtime          = "provided.al2"
  source_code_hash = filebase64sha256("../lambda.zip")
  timeout          = 15
}

# IAM
resource "aws_iam_role" "lambda_role" {
  name               = "shoppinglistapi-lambda-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cloudwatch-attach" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "s3-policy" {
  name        = "shoppinglistapi-lambda-s3-policy"
  description = "Grants lambda access to s3"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": "arn:aws:s3:::*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ssm-policy-attach" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.ssm-policy.arn
}

resource "aws_iam_policy" "ssm-policy" {
  name        = "shoppinglistapi-lambda-ssm-policy"
  description = "Grants lambda access to ssm"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter"
      ],
      "Resource": "arn:aws:ssm:::*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "s3-policy-attach" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.s3-policy.arn
}

# ALB
resource "aws_lb_target_group" "target" {
  name        = "shoppinglistapi"
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "server" {
  target_group_arn = aws_lb_target_group.target.arn
  target_id        = aws_lambda_alias.server_live.arn
  depends_on       = [aws_lambda_permission.server_live]
}

resource "aws_lb_listener_rule" "server" {
  listener_arn = data.terraform_remote_state.stinkyfingers.outputs.stinkyfingers_https_listener
  priority     = 41
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.target.arn
  }
  condition {
    path_pattern {
      values = ["/shoppinglistapi/*"]
    }
  }
  depends_on = [aws_lb_target_group.target]
}

# db
resource "aws_s3_bucket" "shoppinglistapi" {
  bucket = "shoppinglistapi"
}

resource "aws_s3_bucket_policy" "shoppinglistapi_s3" {
  bucket = "shoppinglistapi"
  policy = data.aws_iam_policy_document.allow_lambda_s3.json
}

data "aws_iam_policy_document" "allow_lambda_s3" {
  statement {
    principals {
      type        = "AWS"
      identifiers = [aws_iam_role.lambda_role.arn]
    }
    actions = [
      "s3:*"
    ]
    resources = [
      "arn:aws:s3:::shoppinglistapi",
      "arn:aws:s3:::shoppinglistapi/*"
    ]
  }
}

resource "aws_dynamodb_table" "sl_users" {
  billing_mode     = "PAY_PER_REQUEST"
  hash_key         = "id"
  name             = "sl_users"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "sl_collections" {
  billing_mode     = "PAY_PER_REQUEST"
  hash_key         = "id"
  name             = "sl_collections"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "sl_lists" {
  billing_mode     = "PAY_PER_REQUEST"
  hash_key         = "id"
  name             = "sl_lists"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "sl_items" {
  billing_mode     = "PAY_PER_REQUEST"
  hash_key         = "id"
  range_key        = "name"
  name             = "sl_items"

  attribute {
    name = "id"
    type = "S"
  }
  attribute {
    name = "name"
    type = "S"
  }
}

data "terraform_remote_state" "shoppinglistapi" {
  backend = "s3"
  config = {
    bucket  = "remotebackend"
    key     = "shoppinglistapi/terraform.tfstate"
    region  = "us-west-1"
    profile = var.profile
  }
}