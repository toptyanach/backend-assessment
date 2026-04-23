# Answers

**Q1: In your implementation, what happens if `source.Withdraw()` succeeds but `dest.Deposit()` fails? Show the exact state of both accounts and the returned plan.**
Because the Usecase simply constructs a `Plan` and does not apply anything to the database directly, the failure of `dest.Deposit()` will immediately return a domain error and abort the process.
- **Plan returned:** `nil` (along with the error).
- **In-memory state:** The `source` object has a reduced balance in memory, and the `dest` object remains unchanged. 
- **Database state:** Completely unchanged. The service layer never receives a plan to commit, and the modified in-memory objects will simply be garbage collected.

**Q2: The buggy code applies mutations one at a time. Why is this a problem? Give a specific failure scenario.**
Applying database mutations sequentially without a single transaction block breaks atomicity. 
**Scenario:** User A transfers $100 to User B. 
1. `uc.db.Apply(mutation1)` executes successfully. User A's balance decreases by $100 in the database.
2. The network drops, the database crashes, or a validation error occurs exactly before `uc.db.Apply(mutation2)` is executed.
**Result:** User A lost $100, but User B never received it. The system is in an inconsistent state, and $100 has permanently vanished from the system.

**Q3: Your `UpdateMut` should only include dirty fields. If an account has `balance` changed but `status` unchanged, the mutation should NOT include `status`. Why does this matter for concurrent updates?**
Including only dirty fields prevents "lost updates" in concurrent environments. If we update all fields (including unchanged ones), we risk overwriting valid data that was just modified by another concurrent transaction.
**Example:** Transaction 1 modifies the `balance`. At the exact same time, Transaction 2 (e.g., an admin action) changes the `status` of that same account to 'LOCKED'. If Transaction 1's mutation includes the old, unchanged `status`, when it saves, it will overwrite Transaction 2's 'LOCKED' status back to 'ACTIVE', destroying data. Dirty fields ensure we only touch what we deliberately changed in our specific transaction.

**Q4: Look at this alternative approach ... What problem does this cause that the dirty-field approach avoids?**
This alternative "always include all fields" approach causes the exact concurrency issue described in Q3. By returning every single field regardless of whether it was modified, it creates massive data races. Transactions will blindly overwrite columns they had no intention of modifying, leading to data corruption and lost updates when multiple processes interact with the same entity simultaneously.