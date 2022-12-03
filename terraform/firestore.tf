resource "google_firestore_index" "active_holes" {
  collection = "bg_${local.env}_Hole"

  fields {
    field_path = "Active"
    order      = "ASCENDING"
  }

  fields {
    field_path = "CreatedAt"
    order      = "DESCENDING"
  }
}

resource "google_firestore_index" "user_best_hole_submission" {
  # collection = "bytegolf_UserBestHoleSub_${local.env}"
  collection = "bg_${local.env}_Submission"

  fields {
    field_path = "BGID"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Correct"
    order      = "ASCENDING"
  }

  fields {
    field_path = "HoleID"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Length"
    order      = "ASCENDING"
  }

  # Use the older submission has the user's best
  fields {
    field_path = "SubmittedTime"
    order      = "ASCENDING"
  }
}

resource "google_firestore_index" "user_hole_submissions" {
  # collection = "bytegolf_UserHoleSubs_${local.env}"
  collection = "bg_${local.env}_Submission"

  fields {
    field_path = "BGID"
    order      = "ASCENDING"
  }

  fields {
    field_path = "HoleID"
    order      = "ASCENDING"
  }

  # Get most recent first
  fields {
    field_path = "SubmittedTime"
    order      = "DESCENDING"
  }
}

resource "google_firestore_index" "user_submissions" {
  # collection = "bytegolf_UserSubs_${local.env}"
  collection = "bg_${local.env}_Submission"

  fields {
    field_path = "BGID"
    order      = "ASCENDING"
  }

  # Get most recent first
  fields {
    field_path = "SubmittedTime"
    order      = "DESCENDING"
  }
}

resource "google_firestore_index" "best_hole_submissions_lang" {
  collection = "bg_${local.env}_Submission"

  fields {
    field_path = "Correct"
    order      = "ASCENDING"
  }

  fields {
    field_path = "HoleID"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Language"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Length"
    order      = "ASCENDING"
  }

  fields {
    field_path = "SubmittedTime"
    order      = "ASCENDING"
  }

  depends_on = [
    google_firestore_index.active_test_cases
  ]
}

resource "google_firestore_index" "best_hole_submissions" {
  collection = "bg_${local.env}_Submission"

  fields {
    field_path = "Correct"
    order      = "ASCENDING"
  }

  fields {
    field_path = "HoleID"
    order      = "ASCENDING"
  }

  fields {
    field_path = "Length"
    order      = "ASCENDING"
  }

  fields {
    field_path = "SubmittedTime"
    order      = "ASCENDING"
  }
}

resource "google_firestore_index" "active_test_cases" {
  collection = "bg_${local.env}_Test"

  fields {
    field_path = "Active"
    order      = "ASCENDING"
  }

  fields {
    field_path = "CreatedAt"
    order      = "DESCENDING"
  }
}