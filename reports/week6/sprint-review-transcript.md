# Sprint 4 Review + Customer Trial / UAT Transcript (10 July 2026)

- **Meeting:** Sprint 4 Review combined with the Week 6 customer trial /
  customer-directed UAT, the customer-facing documentation review, and the
  transition-readiness discussion (one recorded session covers all of these;
  the private recording link and exact private timecodes are submitted
  through Moodle).
- **Participants:** Mikhail (PM / QA / Backend), Nurislam (QA / Backend),
  Sanzhar (Scrum Master / QA), and the Customer. Sergei (Scrum Master / QA)
  could not attend this session; his feature slot was covered by the
  teammates.
- **Permissions:** The Customer permitted recording before it started and
  permitted publication of a public transcript.
- **Sanitization note:** This transcript is translated into English,
  shortened, and cleaned for readability without changing meaning. The
  customer's name is replaced with "Customer". Long screen-share
  demonstration segments are condensed; timestamps are approximate
  minutes:seconds into the recording.

---

**00:00**
Mikhail: Hello, and thank you for joining us. Today we want to show you the Week 6 trial version with the remaining features — bulk accrual, the lots audit API, and the multi-key autotester you asked for — and then discuss the documentation and the transition. As always: may we record the meeting, and may we publish a cleaned transcript afterwards?

**00:40**
Customer: Greetings. Yes, you may record, and yes, you may publish the transcript. Sure, let's begin.

**01:12**
Mikhail: Sergei could not join us today, so we will cover his part between us. Nurislam, please start.

**01:30**
Nurislam: Hello. This is the autotester extension you requested last week — running parallel requests with multiple distinct idempotency keys. On the Autotester tab there should be a test-mode selector... *(opens the deployed instance; the selector is not present in the served page)* — one moment, the deployed version is not showing the new controls.

**04:55**
Customer: I see. So the feature is there but the deployment is not right?

**05:10**
Nurislam: Yes — the check itself is implemented in the engine and passes in our test suite; the page you see is serving a stale interface. *(shows the multi-key check running from the console tool instead; demonstration continues, Customer directing)*

**09:20**
Mikhail: While we sort that out, let me show the lots audit for support work. *(demonstrates `GET /v1/users/{id}/lots`: the paginated envelope, the `status=active|expired|exhausted` filter, an invalid filter returning 400; Customer directing)*

**14:45**
Customer: Okay, this part works as I would expect. Good.

**15:30**
Sanzhar: Now the bulk accrual for campaigns. In the Accrue Points tab there is supposed to be a bulk card with a row editor... *(the card is likewise missing from the deployed page; shows the `POST /v1/accruals/batch` endpoint returning per-item results via Swagger instead)*

**19:05**
Customer: Understood. The functionality underneath seems fine, but I cannot try these two from the interface today. Please fix all of this by next week — the interface must show what the product can do.

**19:40**
Mikhail: Agreed — it is a deployment/embedding problem on our side, we have it tracked already and will fix and redeploy today.

**21:10**
Mikhail: Could we also ask you about the documentation set we sent — the README, the handover document, the setup and deployment instructions?

**23:00**
Customer: I read through it. Overall the documentation is complete — the READMEs cover what I need. One thing I want stated explicitly in the technical documentation, the architecture part: whether this service can be scaled horizontally or not. Write it down either way, with the reasoning.

**24:30**
Mikhail: Noted — we will analyse it and state it explicitly in the architecture documentation. And regarding the transition: are you already using the product, or planning to run it on your side?

**26:15**
Customer: Not yet, and honestly not from your VM. This project was designed for a large company context — a simple deployment on a university machine with weak security is not how we would run it. I cannot tell you the internal structure of our deployment. Once we receive the project, a group of interns will look at it and do the deployment on our side.

**28:40**
Mikhail: Understood. Then what do you expect from us for the final delivery next week?

**29:10**
Customer: Overall you will simply show me that all of the main features work correctly and that the documentation is written properly. Check the horizontal-scaling question, fix the interface issue, and polish the project — then I consider the delivery complete. After the handover I will take care of the project on my own.

**31:00**
Mikhail: Great. So for next week: fix the UI issue, assess and document horizontal scaling, and complete the delivery. Thank you for your time.

**31:30**
Customer: Thank you all. See you next week.
