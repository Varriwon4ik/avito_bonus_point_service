# Customer Review Transcript — Sprint 2 (Assignment 4)

> Cleaned English transcript of the recorded Sprint Review and customer-executed
> UAT session held on 28 June 2026. It has been cleaned for readability without
> changing meaning; timestamps are placed on separate lines. The Customer's
> identity and any confidential business details are redacted; the customer is
> referred to as "Customer". Recording permission and permission to publish this
> transcript were granted at the start of the meeting. The private recording link
> and exact recording timecodes are provided through Moodle.

**Participants:** Mikhail (Project Manager / QA / Backend), Nurislam (QA /
Backend), Sanzhar (Scrum Master / QA), Sergey (Scrum Master / QA), Customer.

---

[00:00]

Sanzhar (Scrum Master): Thanks for joining. Before we start — are you OK with us
recording this session, and with us publishing a cleaned, sanitized transcript in
the public repository?

Customer: Yes, that's fine for both.

[00:00] — Sprint Review begins

Mikhail (PM): This Sprint our goal was quality and automation rather than new
end-user features. We delivered three items: the CI pipeline, the autotester, and
pagination for transaction history. I'll walk through everything we committed,
file by file.

[02:10]

Mikhail: Starting with the CI pipeline (US-14) — every push and pull request now
builds the service, runs the full test suite under the race detector, and blocks
merging if anything fails.

Customer: So nothing reaches the main branch without passing these checks?

Mikhail: Correct. The default branch is protected and requires the checks to
pass and a reviewer to approve.

[09:40]

Sanzhar: Next, the autotester (US-15). It lets us define a test scenario —
amount, TTL, how many parallel requests — store it, and replay it against a
running instance, including concurrent requests. It then reports the resulting
balance and ledger so we can confirm the behaviour.

Customer: That's close to what I asked for last time — I wanted to be able to see
that the new behaviour holds up against the older version's expectations.

[16:05]

Nurislam: Pagination (US-09): the transaction history endpoint now takes `page`
and `offset` and returns the page plus the total count, instead of one unbounded
list. Invalid parameters return a clear 400 error.

Customer: Good. Please show me that one live in a moment.

[23:09] — Customer-executed UAT begins

Sanzhar: For acceptance testing we prepared three options: deliver the product to
you to test independently; a remote screen-shared demo where you direct us; or a
full walkthrough of the code and behaviour where you inspect each decision we
made. Which would you prefer?

Customer: Let's do the third — I want to go through the code and the decisions
with you.

Sergey: Sounds good. Let's start with reading a balance [UAT-001], then the
hold/confirm/cancel redemption flow [UAT-002], then the paginated history
[UAT-003].

[24:30]

Customer: Show me the balance read and the query behind it.

Nurislam: Here is the handler and the query — available, held, total, and points
expiring soon.

Customer: That matches what I expected. [redacted — internal pricing example]

[31:50]

Customer: Now the redemption flow.

Mikhail: A hold reserves the points; confirm finalizes the spend; cancel returns
them. Concurrent debits on the same user are serialized in the database, so we
never overspend.

Customer: And that's the part you proved with the automated concurrency test?

Mikhail: Yes — it's one of our quality requirement tests now, run on every change.

[38:15]

Customer: Finally, the paginated history.

Nurislam: Page one, page two, and an invalid offset returning a 400.

Customer: That's exactly what I was after.

[42:40]

Customer: Overall I'm satisfied with this version and with the decisions you've
made. My main ask stays the same: keep building out the automated tests that
demonstrate the changes are valid against the earlier behaviour.

Sanzhar: Understood — we'll extend regression coverage over the older code paths
next Sprint.

[44:05]

Sanzhar: Thank you. We'll send the summary and the action points.

[end]
