---
name: Story issue template
about: Starting point for new stories
title: 'clear, concise summary'
labels: ''
assignees: ''
---

# At a glance

[comment]: # "Begin with a short summary so intent can be understood at a glance."

[comment]: # "In order to: some objective or value to be achieved"
[comment]: # "as a: stakeholder"
[comment]: # "I want: some new feature"

**In order to** 
**as a**
**I want**

# Acceptance Criteria

We use [DRY](https://docs.behat.org/en/latest/user_guide/writing_scenarios.html#backgrounds) [behavior-driven development](https://en.wikipedia.org/wiki/Behavior-driven_development#Behavioral_specifications) wherever possible.

[comment]: # "ACs should be clearly demoable/verifiable whenever possible."
[comment]: # "Given: the initial context at the beginning of the scenario"
[comment]: # "when: the event that triggers the scenario"
[comment]: # "then: the expected outcome(s)"
[comment]: # "Repeat scenarios as needed, or repeat behaviors and lists within a scenario as needed."

[comment]: # "The scenario should be a short, plain language description."
[comment]: # "Feeling repetitive? Apply the DRY (Don't Repeat Yourself) principle!"

### Scenario: 

**Given**   
**when**  
...

[comment]: # "Each task should be a verifiable outcome"
```[tasklist]
### then... 
- [ ] [a thing happens]
```

### Shepherd

[comment]: # "@ mention shepherds as we move across the board."

* UX shepherd: 
* Design shepherd: 
* Engineering shepherd: 

# Background

[comment]: # "Any helpful contextual notes or links to artifacts/evidence, if needed"

# Security Considerations

Required per [CM-4](https://nvd.nist.gov/800-53/Rev4/control/CM-4).

[comment]: # "Our SSP says 'The team ensures security implications are considered as part of the agile requirements refinement process by including a section in the issue template used as a basis for new work.'"
[comment]: # "Please do not remove this section without care."
[comment]: # "Note any security concerns that might be implicated in the change. 'None' is OK, but we must be explicit here."

---

<details>
  <summary>Process checklist</summary>
  
- [ ] Has a clear story statement
- [ ] Can reasonably be done in a few days (otherwise, split this up!)
- [ ] Shepherds have been identified
- [ ] UX youexes all the things
- [ ] Design designs all the things
- [ ] Engineering engineers all the things
- [ ] Meets acceptance criteria
- [ ] Meets [QASP conditions](https://derisking-guide.18f.gov/qasp/)
- [ ] Presented in a review
- [ ] Includes screenshots or references to artifacts
- [ ] Tagged with the sprint where it was finished
- [ ] Archived

### If there's UI...
- [ ] Screen reader - Listen to the experience with a screen reader extension, ensure the information presented in order
- [ ] Keyboard navigation - Run through acceptance criteria with keyboard tabs, ensure it works. 
- [ ] Text scaling - Adjust viewport to 1280 pixels wide and zoom to 200%, ensure everything renders as expected. Document 400% zoom issues with USWDS if appropriate.


</details>