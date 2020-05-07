resource "aws_lambda_function" "default" {
  filename         = "${var.project}.zip"
  function_name    = var.project
  role             = aws_iam_role.lambda_iam_role.arn
  handler          = "lambda"
  source_code_hash = filebase64sha256("${var.project}.zip")
  runtime          = "go1.x"
}

resource "aws_lambda_alias" "default" {
  name             = var.project
  description      = "Use latest version as default"
  function_name    = aws_lambda_function.default.function_name
  function_version = "$LATEST"
}

resource "aws_lambda_permission" "ses" {
  statement_id   = "AllowExecutionFromSES"
  action         = "lambda:InvokeFunction"
  function_name  = aws_lambda_function.default.function_name
  principal      = "ses.amazonaws.com"
  source_account = data.aws_caller_identity.current.account_id
}
