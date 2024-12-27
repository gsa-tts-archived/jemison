---
toc_md_description: toc-md-value
---


[TOC]

## how we work

* **Warmth**. We are, first and foremost, human beings worthy of respect and compassion, generally working under challenging conditions. We all experience our day-to-day differently. While our cup may feel more or less full on any given day, we aspire to bring our most compassionate selves to the work we do and our interactions with our colleagues.
* **Trust**. Trust is earned through words and actions. It comes from doing what we say we will do, when we say we will do it. It comes from saying when we misjudged, and our goals might not be achieved the way we hoped. Open, honest communication is foundational to trust.
    * We value accountability, transparency, and regular open/frank communication. These are all things that lay foundations for *trust*.
* **Working in the open**. As federal employees, our work grows the public domain. As much as possible, we also hope our words and actions might model how this work can also take place openly. In doing so, it encourages others to contribute to our efforts of improving the government and it’s processes *for the benefit of the people*.
* **Collaboration**. We value working together. It is easier to brainstorm, to explore ideas, to get unstuck, and to push each other when we are working together.
* **Documentation**. Our work should be documented. More specifically, we should 1) try and document the rationale for important design decisions (without which, we might “roll back” a decision that was made for good reason), 2) document work products (so that others can pick them up and continue the work), and 3) document reflections (so that we can communicate what and why about our work led to success or failure).
* **Be bold**. In open, collaborative work, being bold means many things. It might mean not asking permission, especially when using tools where changes can easily be rolled back. With regards to our thinking, it means avoiding mental locks, like believing we have to do things that are logical, or that we have to come up with the right answer. While it is true that we have milestones to achieve, part of good agile practice is also knowing how and when to pivot, and that occasionally requires us to be bold. 

### mental locks

Roger von Oech wrote the book A Whack on the Side of the Head. It’s a fun book. A light read.
At the core of this (small) book are ten mental locks that we tend to constrain ourselves with. "Is that the right answer?" "Are we following the rules?" "Can we be practical for a moment?"

These are von Oech's mental locks. Be conscious of when you are using them to keep yourself from exploring ideas. They are pervasive in government.

1. The right answer
2. That's not logical
3. Follow the rules
4. Be practical
5. Play is frivolous
6. That's not my area
7. Don't be foolish
8. Avoid ambiguity
9. To err is wrong
10. I'm not creative




## product principles

* **We work in English and Spanish**. We will expand this list, but these are our primary target languages to start. 
* **We work with our partners**. We are a search engine by and for government. Understanding their needs is critical to our success.
* **We work with our users**. Our partners are agencies, our users are the public. If the public cannot find what they need, we have work to do.

## technical principles

* We aspire to systemic simplicity.
* We do not optimize prematurely.
* We prefer simpler, fast-to-market solutions we can iterate on over larger, slower, more complex solutions that might work no better.
* **Service-oriented**. Jemison is a service-oriented architecture: if a single service seems to be getting to complex, perhaps it should be two?
* **Use the queue**. Services should communicate via the queue.
* **Automate**. Everything, where possible, should be automated.
* **Non-interactive**. To start, we drive everything via configuration files. Changes require a re-deployment. This is testable.
* **We deploy statelessly**. As much as possible, our deploys are stateless TF deploys. 
