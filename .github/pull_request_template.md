We prefer smaller PRs that merge more often.

## Haiku-length summary



### Additional details 



## PR Checklist: Submitter

- [ ]   Link to an issue if possible. If there’s no issue, describe what your branch does. Even if there is an issue, a brief description in the PR is still useful.
- [ ]   List any special steps reviewers have to follow to test the PR. For example, adding a local environment variable, creating a local test file, etc.
- [ ]   For extra credit, submit a screen recording like [this one](https://github.com/GSA-TTS/FAC/pull/1821).
- [ ]   Make sure you’ve merged `main` into your branch shortly before creating the PR. (You should also be merging `main` into your branch regularly during development.)
- [ ]   Make sure that whatever feature you’re adding has tests that cover the feature. This includes test coverage to make sure that the previous workflow still works, if applicable.
- [ ]   Make sure the E2E tests pass.
- [ ]   Do manual testing locally. 
  - [ ]   If that’s not applicable for some reason, check this box.
- [ ]   Once a PR is merged, keep an eye on it until it’s deployed to dev, and do enough testing on dev to verify that it deployed successfully, the feature works as expected, and the happy path for the broad feature area still works.

## PR Checklist: Reviewer

- [ ]   Pull the branch to your local environment and run `make macup ; make e2e"` (FIXME)
- [ ]   Manually test out the changes locally
  - [ ]   Check this box if not applicable in this case.
- [ ]   Check that the PR has appropriate tests.

The larger the PR, the stricter we should be about these points.
