resource aws_lambda_function lambda_function {
  filename         = format("%s.zip", var.project)
  function_name    = var.project
  role             = aws_iam_role.lambda_iam_role.arn
  handler          = "lambda"
  source_code_hash = filebase64sha256(format("%s.zip", var.project))
  runtime          = "go1.x"
  depends_on       = [aws_iam_role_policy_attachment.lambda, aws_cloudwatch_log_group.lambda_log_group]
}

resource aws_lambda_alias lambda_alias {
  name             = var.project
  description      = "Use latest version as default"
  function_name    = aws_lambda_function.lambda_function.function_name
  function_version = "$LATEST"
}

resource aws_lambda_permission lambda_ses_permission {
  statement_id   = "AllowExecutionFromSES"
  action         = "lambda:InvokeFunction"
  function_name  = aws_lambda_function.lambda_function.function_name
  principal      = "ses.amazonaws.com"
  source_account = data.aws_caller_identity.current.account_id
}
