resource "google_pubsub_topic" "bgcompiler" {
  name = "bgcompiler-${local.env}"
  message_retention_duration = "18000s"
}

resource "google_pubsub_subscription" "bgcompiler_sub" {
  # TODO: Add env to name.
  name = "bgcompiler-sub"
  topic = google_pubsub_topic.bgcompiler.name
  ack_deadline_seconds = 30
  retain_acked_messages = false
  enable_exactly_once_delivery = true
}