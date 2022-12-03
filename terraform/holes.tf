# resource "google_firestore_document" "hello_world" {
#   project = local.project
#   collection = google_firestore_index.active_holes.collection
#   document_id = "testing_hello_world"
#   fields = "{\"Active\": {\"booleanValue\": true},\"Difficulty\": {\"stringValue\": \"EASY\"}, \"ID\":{\"stringValue\": \"testing_hello_world\"},\"Name\":{\"stringValue\": \"Testing Hello, World!\"},\"Question\": {\"stringValue\": \"Print \\\"Hello, World!\\\" to console.\"}}"
# }