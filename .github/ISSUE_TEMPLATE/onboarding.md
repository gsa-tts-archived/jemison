---
name: 🥳 Onboarding
about: An onboarding checklist for folks joining the team!
title: Onboarding checklist for {GitHub Username}
labels: 'backlog'
assignees: ''
---

**Instructions for the issue creator:** 
1. Fill in the GitHub username of the new person (if known) and the name of the onboarding buddy. 
2. Remove the horizontal line below and everything above it (these instructions).
3. Remove sections that do not apply (based on role)

---

# Welcome to the team @{GitHub Username}!

{ONBOARDING BUDDY'S GH USERNAME} will be your onboarding buddy and can help you if you get stuck on any of the steps in the checklist below. There's also a separate checklist that your onboarding buddy will tackle!

## For all team members

All team members will want to take the following steps:

- [ ] Review [the project README](https://github.com/GSA-TTS/jemison/blob/main/docs/README.md)
- [ ] Ask in `#search-general` for team members to reach out and setup an intro coffee
- [ ] Complete OLU trainings
- [ ] Review the TTS handbook
  - [ ] https://handbook.tts.gsa.gov/getting-started/
  - [ ] https://handbook.tts.gsa.gov/about-us/tts-history/
  - [ ] Skim through: https://handbook.tts.gsa.gov/
- [ ] Review our [agile processes](https://github.com/GSA-TTS/jemison/blob/main/docs/process/agile.md)
  - [ ] If you are unfamiliar with agile, [this is a good overview](https://digital.gov/event/2019/11/04/foundations-agile-i/)
- [ ] Familiarized yourself with the [Jemison architecture](https://github.com/GSA-TTS/jemison/blob/main/docs/architecture/index.md)

### Systems access

- [ ] Join/get added to the team Slack channels
- [ ] Get added to the team shared drive
- [ ] (Feds) get added to the Federal shared drive
- [ ] Have the team calendar owner share the calendar invite link to our new team member
- [ ] Get added to our ZD helpdesk instance
- [ ] Get access to our GitHub repositories
  - [ ] If you're not already a member of BOTH the [GSA](https://github.com/orgs/GSA/people) and [GSA-TTS](https://github.com/orgs/GSA-TTS/people) GitHub organizations, [follow the process outlined here](https://github.com/GSA/GitHub-Administration#joining-the-gsa-enterprise-organization), requesting access to the appropriate org(s). When it's time to send the mail, look up the Product Lead (see [the staffing list](https://docs.google.com/document/d/1g8nYqYS_ifFlZB-DBgfeSoJRMB__EqWsmLnacyk-bDI/edit#heading=h.us8xylqg455c)) then start with this template, mailing from your GSA email address:

    ```text
    To: gsa-github.support@gsa.gov
    Cc: [the search lead]
    Subject: GSA and GSA-TTS GitHub organization membership
    Body:
    [Edit the following text to put in your GitHub username and your GSA email address, then remove this instruction text.]
    Please add my GitHub account (https://github.com/myusername) to the following GitHub organizations:
    - https://github.com/GSA
    - https://github.com/GSA-TTS
    
    I will be working on the Search.gov project:
    - https://github.com/GSA-TTS/jemison
    
    I have cc'd the search.gov lead for awareness.
        
    Thank you!
    ```        
    - (Note this step could take a few days; humans handle these requests.)
  - [ ] Once you are added to the GSA-TTS org, ask [the person(s) with the "Maintainer" role to add you to the `searchgov-team` team](https://github.com/orgs/GSA-TTS/teams/searchgov-team/members). This will grant you read/write access to our repositories.

## For UX/design

- [ ] Review the [design onboarding document](TBD)
- [ ] Visit https://touchpoints.digital.gov/ and set up a touchpoints account
- [ ] If you don't already have a FigJam license, ping in the general channel that you need one. 

## For CX/content

TBD

## For Engineering

- [ ] Familiarize yourself with Go, Jsonnet, and Cloud.gov—tools used in this project.
  - [ ] If it's not already set up on your machine & account, enable [commit signing](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits). Also see the pinned messages in the Slack development channel.
  - [ ] If you're not already, get [setup with Cloud.gov](https://cloud.gov/docs/getting-started/setup/)
    - [ ] Once your account exists, make a PR to [add yourself to the list of developers](TBD) with access to our spaces.
  - [ ] Practice deploying a [python application](https://github.com/cloud-gov/cf-hello-worlds/tree/main/python-flask) to Cloud.gov using the Cloud.gov command line interface (CLI): https://cloud.gov/docs/getting-started/your-first-deploy/.
- [ ] Get set up for [local development](https://github.com/GSA-TTS/jemison/docs/process/development.md) so you can start contributing
- [ ] Review the main body (~60p) of [NIST SP 800-161 Rev.1](https://csrc.nist.gov/pubs/sp/800/161/r1/upd1/final) on supply chain risk management. Reply to this onboarding ticket in a comment when you have completed this review.


## For your onboarding buddy

Note: If you're not able to do any of these yourself, feel free to leverage your onboarding buddy (or any number of team members) to help make sure things happen.
