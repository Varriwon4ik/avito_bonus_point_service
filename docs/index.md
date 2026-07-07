# Bonus Points Ledger Service — Documentation

Maintained documentation for the **Bonus Points Ledger Service**, a REST-like
service managing an online store's bonus-points program: configurable point
expiry with validated TTL bounds, transactional balance mutations with
row-level locking, two-phase debits, idempotent operations, FIFO-by-expiry
consumption, transaction labels, paginated history, observability, a web UI
with a built-in autotester, and Swagger/OpenAPI docs.

- **Repository:** <https://github.com/Varriwon4ik/avito_bonus_point_service>
- **Setup and API reference:** the
  [root README](https://github.com/Varriwon4ik/avito_bonus_point_service#readme)

## Documentation map

| Area | Page |
|---|---|
| System structure, views, and diagrams | [Architecture](architecture/README.md) |
| Why key decisions were made | [ADR index](architecture/README.md#architecture-decision-records-adr-index) |
| How the team works and manages configuration | [Development process](development-process.md) |
| Sprint-by-Sprint delivery plan | [Roadmap](roadmap.md) |
| Completion standard | [Definition of Done](definition-of-done.md) |
| Testing status, coverage, CI gates | [Testing](testing.md) |
| Measurable non-functional requirements | [Quality requirements](quality-requirements.md) |
| Automated tests verifying the QRs | [Quality requirement tests](quality-requirement-tests.md) |
| Customer-executable acceptance scenarios | [User acceptance tests](user-acceptance-tests.md) |
| Stable user-story registry | [User stories](user-stories.md) |

This site is built from the repository's `docs/` directory with MkDocs and
republished automatically on every merge to `main`.
