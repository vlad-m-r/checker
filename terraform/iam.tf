data aws_iam_policy_document assume {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = [
      "sts:AssumeRole"
    ]
  }
}

resource aws_iam_role lambda_iam_role {
  name               = var.project
  assume_role_policy = data.aws_iam_policy_document.assume.json
}

data aws_iam_policy_document lambda {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]

    resources = ["*"]
  }

  statement {
    effect = "Allow"

    actions = [
      "ses:SendEmail",
      "ses:SendRawEmail"
    ]

    resources = ["*"]
  }
}

resource aws_iam_policy lambda_iam_policy {
  name        = var.project
  description = "Allows to sent emails with SES and store logs in Cloudwatch"
  policy      = data.aws_iam_policy_document.lambda.json
}

resource aws_iam_role_policy_attachment lambda {
  role       = aws_iam_role.lambda_iam_role.name
  policy_arn = aws_iam_policy.lambda_iam_policy.arn
}