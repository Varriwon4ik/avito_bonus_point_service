# Sprint 3 Review + UAT Transcript (3 July 2026)

- **Meeting:** Sprint 3 Review combined with customer-directed UAT
  (one recorded session covers both; the private recording link and exact
  private timecodes are submitted through Moodle).
- **Participants:** Mikhail (PM / QA / Backend), Nurislam (QA / Backend),
  Sanzhar (Scrum Master / QA), Sergei (Scrum Master / QA), and the Customer.
- **Permissions:** The Customer permitted recording before it started and
  permitted publication of a public transcript.
- **Sanitization note:** This transcript is translated into English, shortened,
  and cleaned for readability without changing meaning. The customer's name is
  replaced with "Customer". Long screen-share demonstration segments (most of
  the meeting) are condensed; timestamps are approximate minutes:seconds into
  the recording.

---

**00:00**
Mikhail: Thank you for joining us today. We are almost done with the project, and today we would like to show you the UI implementation of the features we proposed last week. May we record the meeting?

**01:05**
Customer: Greetings — yes, you may record our meeting. I see your scope for this week, so please, let's begin.

**01:46**
Mikhail: Great. I would like to ask my teammates to present their work, and I want Nurislam to begin, as we are currently having some troubles with our repository. I am afraid we will have to show you the undeployed version.

**02:07**
Customer: It's okay. Please, let's start.

**02:13**
Nurislam: Hello. This is the frontend implementation of HTTP responses. As you can see, I have an "Accrue points" window, and once I add the points, this green text appears with the code 200, which means everything is OK. Next I will try to debit points from a non-existing user, and here you can see the red code 404 — the user was not found. *(demonstration continues, Customer directing)*

**12:07**
Customer: Okay, I see. But you could change how the response codes are shown here. It's not critical — just if you have the courage.

**13:21**
Nurislam: Okay, I understand.

**13:25**
Customer: Yeah — then I'd have to tell you why that approach is more convenient, right.

**15:48**
Customer: Okay, do we have something else to inspect? Because I would honestly like to see your plans for the next Sprint.

**16:31**
Mikhail: Sure, let me list the user stories that are left. It's US-01 and US-02.

**18:10**
Customer: Mhm, yes, good. But let me tell you my other point: you have to all work as a team — as a whole.

**20:15**
Mikhail: Yeah, we are trying to work as a whole. We just have a few problems that come up, and because of them we may seem to work solo — others are sometimes unaware of how to complete a certain task, and walking them through it would take a lot of time, so we discuss it and then let them watch the process.

**22:44**
Mikhail: So, I would like to show my autotester UI implementation. *(demonstration of the Autotester tab, Customer directing)*

**24:42**
Customer: Alright — does it work on a single idempotency key or on multiple ones?

**24:50**
Mikhail: It works only on a single one.

**25:01**
Customer: Oh, I see now. I'd say it is good, but it would be nice if there were a feature to run a test with multiple idempotency keys in parallel requests and check whether that works.

**27:45**
Mikhail: Sure, that is a great idea — I've added it to our backlog now, thank you. Now let Sanzhar present.

**28:15**
Sanzhar: Hello. This is my implementation of the labels feature. You can add one from a preset list, or add a custom one yourself. *(demonstration of preset and custom labels, Customer directing)*

**33:07**
Customer: Right, good — the custom labels feature is specifically good, because you can assign a certain product a label on its transaction and have a convenient tool to confirm debits. Good, yes.

**38:13**
Mikhail: Okay, and now I would like Sergei to demonstrate his feature.

**38:20**
Sergei: Right — this is the TTL validation for accruals. *(demonstration of the configured TTL bounds and the 400 response outside the range)*

**39:11**
Customer: Good, good. I see the changes — works great. Okay, we have discussed your plans too now and seen the current version. But please don't forget to work as a team, guys — this is essential in your work.

---

*End of sanitized transcript. The backlog item created during the meeting is
[US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50).*
