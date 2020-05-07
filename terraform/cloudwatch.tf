resource "aws_cloudwatch_event_rule" "event_rule" {
  schedule_expression = var.interval
}

resource "aws_cloudwatch_event_target" "check_at_interval" {
  rule = aws_cloudwatch_event_rule.event_rule.name
  arn = aws_lambda_function.default.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_lambda" {
  statement_id = "AllowExecutionFromCloudWatch"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.default.function_name
  principal = "events.amazonaws.com"
  source_arn = aws_cloudwatch_event_rule.event_rule.arn
}