ALTER TABLE pull_request
ADD COLUMN (
  additions int,
  deletions int,
  review_duration bigint
)