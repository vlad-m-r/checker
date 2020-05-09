resource aws_cloudwatch_event_rule event_rule {
  schedule_expression = var.interval
}

resource aws_cloudwatch_event_target check_at_interval {
  rule = aws_cloudwatch_event_rule.event_rule.name
  arn  = aws_lambda_function.lambda_function.arn
}

resource aws_lambda_permission allow_cloudwatch_to_call_lambda {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.event_rule.arn
}


resource aws_cloudwatch_log_group lambda_log_group {
  name              = format("/aws/lambda/%s", var.project)
  retention_in_days = 3
}
