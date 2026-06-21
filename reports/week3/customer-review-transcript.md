# Customer Review Transcript — Sprint 1 (MVP v1)

**Date:** 21 Jun 2026
**Event:** Sprint Review with the customer
**Participants (roles only):** Product Owner / Scrum Master (Mikhail), Developer
(Nurislam), QA Engineer (Sergey), QA Engineer (Sanzhar), Customer (Alexander).

> This transcript has been cleaned for readability without changing meaning.
> Personally identifying and confidential business details have been removed or
> replaced with `[redacted]`. `[inaudible]` marks unclear audio. Recording began
> partway through the meeting; see `customer-review-notes.md` is not required as
> recording and review notes are captured here and in the summary.

---

00:10
Mikhail: So, our team has prepared a proposed MVP v1 scope for our project. It consists of the main features you assigned us last week.

00:21
Alexander (customer): Yes, you can present now.

00:25
Mikhail: Great. First I would like to show the completed US-11, where our team had to implement parallel request handling. Nurislam, please demonstrate your work.

00:40
Nurislam: Yes, so this is the implementation of the solution for the "last-write-wins" issue.

00:45 – 09:26
Nurislam demonstrates the implementation of US-11 (concurrent idempotent-key deduplication and race-safe balance mutation). [Discussion of technical details with the customer during this interval.]

09:26
Alexander: Great, but it would be even better if you would apply similar tests to the older versions to show that our solution has an impact on the project.

10:26
Mikhail: We see, thank you Alexander. Now I would like to ask Sergey to demonstrate how our application works with expiring points. Sergey, please show.

11:13
Sergey: Yes, so this is how our system works with points.

11:20 – 13:47
Sergey demonstrates the implementation of US-12 (points removal / expiry system). [Discussion of technical details with the customer during this interval.]

13:47
Alexander: Okay, I see. A good implementation. Though I would want you to not handle expirations apart from transactions. Instead of just marking points as expired and completely removing them from the database, you could create a new transaction that removes the points and is recorded in the database with the expiration-date handling.

15:21
Mikhail: That is a great suggestion, thank you Alexander. Right now I want Sanzhar to present the last feature we had to implement.

15:45
Sanzhar: Yes, so this is the HTTP-response-codes handling that I implemented.

15:52 – 17:12
Sanzhar demonstrates the implementation of US-13 (HTTP response codes and OpenAPI docs). [Discussion of technical details with the customer during this interval.]

17:12
Alexander: Good, I see good code decisions here.

18:13
Mikhail: So this was our presentation of the currently assigned user stories. Unfortunately we couldn't connect our fifth team member to demonstrate how he implemented US-05. Is it okay if we show it next week?

18:45
Alexander: Yes, no problem.

---

**Note on timestamps:** the interval jumps (e.g. 00:45–09:26) correspond to live
demonstrations during which the customer discussed the technical traits of each
implementation. Only the spoken decision points are transcribed here.

---

**Publication permission:** The customer, Alexander, explicitly permitted public
publication of this sanitized English transcript in the project's public
repository.
