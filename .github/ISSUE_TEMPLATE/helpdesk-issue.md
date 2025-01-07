name: Helpdesk Issue
description: Used as output from helpdesk issue triage.
title: "[HD]: "
labels: ["helpdesk"]
projects: ["GSA-TTS/60"]
assignees:
  - jadudm
body:
  - type: markdown
    attributes:
      value: |
        **Remember to keep all PII out of tickets.**
        
        Do not use names, and do not attach files to this ticket. 
  - type: input
    id: zendesk-link
    attributes:
      label: Zendesk link
      description: Link to the issue in Zendesk
      placeholder: ex. https://fac-something.zendesk.com/something/...
    validations:
      required: true
  - type: checkboxes
    id: fac-components
    attributes:
      label: Search.gov components involved
      description: Select all that apply
      options:
        - label: crawling/indexing
        - label: search results page
        - label: API
        - label: the helpdesk
        - label: static site
        - label: other
    validations:
      required: true
  - type: dropdown
    id: browser
    attributes:
      label: What browser did the user report as using?
      multiple: false
      options:
        - Firefox
        - Chrome
        - Safari
        - Microsoft Edge
        - Other (Opera, Brave, etc.)  
  - type: markdown
    attributes:
      value: |
        Place all files and screenshots in the [Google Drive Helpdesk folder](https://drive.google.com/drive/folders/1L5eCVFlXxMyF98O2us5znZUPYaOx02gI) and link to that folder here.

  - type: input
    id: gdrive-link
    attributes:
      label: Gdrive link
      description: Link to supporting files in GDrive
      placeholder: ex. https://google.com/drive/something/...
    validations:
      required: false  
  - type: textarea
    id: what-happened
    attributes:
      label: What happened?
      description: Summarize the issue the user is experiencing. 
      placeholder: 
    validations:
      required: true